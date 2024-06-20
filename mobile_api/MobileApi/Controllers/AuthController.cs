using System.Net.Mail;
using System.Net;
using Microsoft.AspNetCore.DataProtection;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Options;
using MobileApi.Sources.db.zebra;
using MobileApi.Sources.email;
using System.Text;
using Newtonsoft.Json;

namespace MobileApi.Controllers.Auth {
    [ApiController]
    [Route("auth")]
    public class AuthController : ControllerBase {

        private readonly IAuthService _authService;

        public AuthController(IAuthService authService) {
            _authService = authService;
        }

        [HttpPost]
        [Route("try-sign-in")]
        public async Task<TrySignInRes> TrySignIn(TrySignInReq req) {
            return await _authService.TrySignIn(req);
        }

        [HttpPost]
        [Route("verify-email-code")]
        public async Task<SignInVerifyEmailCodeRes> VerifyEmailCode(SignInVerifyEmailCodeReq req) {
            return await _authService.VerifyEmailCode(req);
        }

        [HttpPost]
        [Route("verify-email-link")]
        public async Task<SignInVerifyEmailCodeRes> VerifyEmailLink(SignInVerifyEmailLinkReq req) {
            return await _authService.VerifyEmailLink(req);
        }

    }

    public interface IAuthService {
        public Task<TrySignInRes> TrySignIn(TrySignInReq req);
        public Task<SignInVerifyEmailCodeRes> VerifyEmailCode(SignInVerifyEmailCodeReq req);
        public Task<SignInVerifyEmailCodeRes> VerifyEmailLink(SignInVerifyEmailLinkReq req);
    }


    public record TrySignInReq(string Email, string DeviceId);
    public record TrySignInRes(TrySignInResStatus Status, DateTime? CodeActiveUntil);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TrySignInResStatus {
        CodeSentToEmail = 1,
        AlreadySentToEmail = 2,
        UserNotExists = 3,
    }


    public record SignInVerifyEmailCodeReq(string Email, string CodeFromEmail);
    public record SignInVerifyEmailLinkReq(string TokenFromRegisterLink, string DeviceId);


    public record SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus Status, string? Email, string? Name, Guid? ClientId);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum SignInVerifyEmailCodeResStatus {
        ValidCode = 1,
        IncorrectCode = 2,
        ToManyAttemps = 3,
        NoSentCode = 4,
        DeviceIdNotMatch = 5,
    }

    public class AuthService : IAuthService {
        private readonly IMemoryCache _memoryCache;
        private readonly ISignInEmailSenderService _emailSenderService;
        private readonly ZebraApplicationContext _db;
        private readonly IDataProtector _dataProtector;

        public AuthService(IMemoryCache memoryCache, ZebraApplicationContext db, ISignInEmailSenderService emailSenderService, IDataProtectionProvider dataProtectionProvider) {
            _memoryCache = memoryCache;
            _db = db;
            _emailSenderService = emailSenderService;
            _dataProtector = dataProtectionProvider.CreateProtector("email_code_encrypt");
        }

        public async Task<TrySignInRes> TrySignIn(TrySignInReq req) {
            var userExists = MobileUserModel.UserExistsWithEmail(req.Email, _db);
            if (!userExists) {
                return new TrySignInRes(TrySignInResStatus.UserNotExists, CodeActiveUntil: null);
            }

            var key = $"auth_{req.Email}";
            var sendStatus = TrySignInResStatus.AlreadySentToEmail;
            var sentCodeInfo = await _memoryCache.GetOrCreateAsync(key, async (memory) => {
                memory.AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(5);

                var code = randomNumbers(4);
                var authByLinkToken = _dataProtector.Protect(new SignInLinkToken(req.Email, code).ToJson());
                await _emailSenderService.SendVerifyCode(req.Email, code, authByLinkToken);
                sendStatus = TrySignInResStatus.CodeSentToEmail;
                return new SentCodeInfo(code, DateTime.Now.AddMinutes(5), req, new VerifyEmailCodeAttempsInfo(Count: 0));
            });

            return new TrySignInRes(sendStatus, sentCodeInfo!.Until);
        }

        public record SignInLinkToken(string Email, string Code) {
            public string ToJson() {
                return JsonConvert.SerializeObject(this);
            }
            public static SignInLinkToken FromJson(string json) => JsonConvert.DeserializeObject<SignInLinkToken>(json)!;
        };
        private static string randomNumbers(int length) {
            var random = new Random();
            var sb = new StringBuilder();
            for (var i = 0; i < length; i++) {
                sb.Append(random.Next(0, 9));
            }

            return sb.ToString();
        }

        public async Task<SignInVerifyEmailCodeRes> VerifyEmailCode(SignInVerifyEmailCodeReq req) {
            var key = $"auth_{req.Email}";
            var sentCodeInfo = _memoryCache.Get<SentCodeInfo>(key);
            if (sentCodeInfo == null) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            var codeExpired = sentCodeInfo.Until < DateTime.Now;
            if (codeExpired) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.AttempsInfo.Count > 5) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.ToManyAttemps, Email: null, Name: null, ClientId: null);
            }

            var isCodeCorrect = sentCodeInfo.Code == req.CodeFromEmail;
            var isGoogleAndAppleTestEmail = sentCodeInfo.User.Email == "maksutov911@gmail.com" && req.CodeFromEmail == "9999";
            if (isCodeCorrect || isGoogleAndAppleTestEmail) {
                var user = await MobileUserModel.GetUserId(_db, sentCodeInfo.User.Email);
                _memoryCache.Remove(key);
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.ValidCode, Email: sentCodeInfo.User.Email, Name: user.Name, ClientId: user.Id);
            } else {
                _memoryCache.Set(key, sentCodeInfo with {
                    AttempsInfo = sentCodeInfo.AttempsInfo with {
                        Count = sentCodeInfo.AttempsInfo.Count + 1
                    }
                }, new MemoryCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10) });

                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.IncorrectCode, Email: null, Name: null, ClientId: null);
            }
        }

        public async Task<SignInVerifyEmailCodeRes> VerifyEmailLink(SignInVerifyEmailLinkReq req) {
            var tokenData = SignInLinkToken.FromJson(_dataProtector.Unprotect(req.TokenFromRegisterLink));

            var key = $"auth_{tokenData.Email}";
            var sentCodeInfo = _memoryCache.Get<SentCodeInfo>(key);
            if (sentCodeInfo == null) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }
            if (req.DeviceId != sentCodeInfo.User.DeviceId) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.DeviceIdNotMatch, Email: null, Name: null, ClientId: null);
            }

            var codeExpired = sentCodeInfo.Until < DateTime.Now;
            if (codeExpired) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.AttempsInfo.Count > 5) {
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.ToManyAttemps, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.Code == tokenData.Code) {
                var user = await MobileUserModel.GetUserId(_db, sentCodeInfo.User.Email);
                _memoryCache.Remove(key);
                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.ValidCode, Email: sentCodeInfo.User.Email, Name: user.Name, ClientId: user.Id);
            } else {
                _memoryCache.Set(key, sentCodeInfo with {
                    AttempsInfo = sentCodeInfo.AttempsInfo with {
                        Count = sentCodeInfo.AttempsInfo.Count + 1
                    }
                }, new MemoryCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10) });

                return new SignInVerifyEmailCodeRes(SignInVerifyEmailCodeResStatus.IncorrectCode, Email: null, Name: null, ClientId: null);
            }
        }

       
    }

    public record VerifyEmailCodeAttempsInfo(int Count);

    public record SentCodeInfo(string Code, DateTime Until, TrySignInReq User, VerifyEmailCodeAttempsInfo AttempsInfo);


    public interface ISignInEmailSenderService {
        public Task SendVerifyCode(string toEmail, string code, string authByLinkToken);
    }

    public class SignInEmailSenderService : ISignInEmailSenderService {
        private readonly EmailSenderConfig _conf;

        public SignInEmailSenderService(IOptions<EmailSenderConfig> conf) {
            _conf = conf.Value;
        }

        public async Task SendVerifyCode(string toEmail, string code, string authByLinkToken) {
            var smtpClient = new SmtpClient("smtp.mail.ru") {
                Port = 587,
                Credentials = new NetworkCredential(_conf.UserName, _conf.UserPwd),
                EnableSsl = true,
            };

            var continueOnMobileLink = $"{_conf.RegisterContinueOnMobileLink}/mobile-app?Type=sign-in&amp;AuthByLinkToken={authByLinkToken}";

            var from = new MailAddress(_conf.UserName, "ZebraCoffee");
            var to = new MailAddress(toEmail);
            var m = new MailMessage(from, to);
            m.Subject = "Авторизация";
            m.Body = EmailMsgs.GetEmailVerifyMsg(code, continueOnMobileLink).Value;
            m.IsBodyHtml = true;

            smtpClient.Send(m);
        }
    }

}
