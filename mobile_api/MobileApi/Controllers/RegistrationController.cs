using System.Net.Mail;
using System.Net;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Caching.Memory;
using MobileApi.Sources.db.zebra;
using System.Text;
using MobileApi.Sources.email;
using Microsoft.Extensions.Options;
using Microsoft.AspNetCore.DataProtection;
using Newtonsoft.Json;

namespace MobileApi.Controllers.Registration {
    [ApiController]
    [Route("registration")]
    public class RegistrationController : ControllerBase {

        private readonly IRegistrationService _registrationService;

        public RegistrationController(IRegistrationService registrationService) {
            _registrationService = registrationService;
        }

        [HttpPost]
        [Route("try-create-user")]
        public async Task<TryCreateUserRes> TryCreateUser(NeedToCreateUserInfo req) {
            return await _registrationService.TryCreateUser(req);
        }

        [HttpPost]
        [Route("create-incognito-user")]
        public async Task<CreateIncognitoUserRes> CreateIncognitoUser(NeedToCreateIncognitoUserInfo req) {
            return await _registrationService.CreateIncognitoUser(req);
        }

        [HttpPost]
        [Route("verify-email-code")]
        public async Task<RegistrateVerifyEmailCodeRes> VerifyEmailCode(RegistrateVerifyEmailCodeReq req) {
            return await _registrationService.VerifyEmailCode(req);
        }
        
        [HttpPost]
        [Route("verify-email-link")]
        public async Task<RegistrateVerifyEmailCodeRes> VerifyEmailLink(RegistrateVerifyEmailLinkReq req) {
            return await _registrationService.VerifyEmailLink(req);
        }

        [HttpPost]
        [Route("try-remove-user")]
        public async Task<TryRemoveUserRes> TryRemoveUser(Guid userId) {
            return await _registrationService.TryRemoveUser(userId);
        }
    }

    public interface IRegistrationService {
        public Task<TryCreateUserRes> TryCreateUser(NeedToCreateUserInfo req);
        public Task<CreateIncognitoUserRes> CreateIncognitoUser(NeedToCreateIncognitoUserInfo req);
        public Task<RegistrateVerifyEmailCodeRes> VerifyEmailCode(RegistrateVerifyEmailCodeReq req);
        public Task<RegistrateVerifyEmailCodeRes> VerifyEmailLink(RegistrateVerifyEmailLinkReq req);
        public Task<TryRemoveUserRes> TryRemoveUser(Guid userId);
    }


    public record NeedToCreateUserInfo(string Email, string ClientName, DateTime? BirthDate, string DeviceId);
    public record NeedToCreateIncognitoUserInfo(string ClientName, string DeviceId);
    public record TryCreateUserRes(TryCreateUserResStatus Status, DateTime? CodeActiveUntil);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryCreateUserResStatus {
        CodeSentToEmail = 1,
        AlreadySentToEmail = 2,
        AlreadyRegistered = 3,
    }

    public record RegistrateVerifyEmailCodeReq(string Email, string CodeFromEmail);
    public record RegistrateVerifyEmailLinkReq( string TokenFromRegisterLink, string DeviceId);
    public record RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus Status, string? Email, string? Name, Guid? ClientId);
    public record CreateIncognitoUserRes(Guid ClientId);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum RegistrateVerifyEmailCodeResStatus {
        ValidCode = 1,
        IncorrectCode = 2,
        ToManyAttemps = 3,
        NoSentCode = 4,
        DeviceIdNotMatch = 5,
    }

    public record TryRemoveUserRes(TryRemoveUserResStatus Status);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryRemoveUserResStatus {
        UserNotExists = 1,
        Ok = 2,
    }

    public class RegistrationService : IRegistrationService {
        private readonly IMemoryCache _memoryCache;
        private readonly IRegistrationEmailSenderService _emailSenderService;
        private readonly ZebraApplicationContext _db;
        private readonly IDataProtector _dataProtector;

        public RegistrationService(IMemoryCache memoryCache, ZebraApplicationContext db, IRegistrationEmailSenderService emailSenderService, IDataProtectionProvider dataProtectionProvider) {
            _memoryCache = memoryCache;
            _db = db;
            _emailSenderService = emailSenderService;
            _dataProtector = dataProtectionProvider.CreateProtector("email_code_encrypt");
        }


        public async Task<CreateIncognitoUserRes> CreateIncognitoUser(NeedToCreateIncognitoUserInfo  req) {

            var key = $"incognito_user_for_{req.DeviceId}";
            var userId = _memoryCache.Get<Guid?>(key);
            if (userId.HasValue) {
                await MobileUserModel.UpdateIncognitoUser(_db, userId.Value, req.ClientName);
            } else {
                userId = await MobileUserModel.InsertIncognitoUser(_db, req.ClientName);
                _memoryCache.Set(key, userId, new MemoryCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(60)});
            }

            return new CreateIncognitoUserRes(userId.Value);
        }

        public async Task<TryCreateUserRes> TryCreateUser(NeedToCreateUserInfo req) {

            var alreadyCreated = _db.Users.Where(x => x.Status == MobileUserModel.UserStatus.Active).Count(x => x.Email == req.Email) > 0;
            if(alreadyCreated) {
                return new TryCreateUserRes(TryCreateUserResStatus.AlreadyRegistered, CodeActiveUntil: null);
            }

            var status = TryCreateUserResStatus.AlreadySentToEmail;
            var key = $"reg_{req.Email}";
            var sentCodeInfo = await _memoryCache.GetOrCreateAsync(key, async (memory) => {
                memory.AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(5);

                var code = randomNumbers(4);
                var authByLinkToken = _dataProtector.Protect(new RegisterLinkToken(req.Email, code).ToJson());
                await _emailSenderService.SendVerifyCode(req.Email, code, authByLinkToken);
                status = TryCreateUserResStatus.CodeSentToEmail;
                return new SentCodeInfo(code, DateTime.Now.AddMinutes(5), req, new VerifyEmailCodeAttempsInfo(Count: 0));
            });

            return new TryCreateUserRes(status, sentCodeInfo!.Until);
        }

        public async Task<RegistrateVerifyEmailCodeRes> VerifyEmailCode(RegistrateVerifyEmailCodeReq req) {
            var key = $"reg_{req.Email}";
            var sentCodeInfo = _memoryCache.Get<SentCodeInfo>(key);
            if(sentCodeInfo == null) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            var codeExpired = sentCodeInfo.Until < DateTime.Now;
            if (codeExpired) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            if(sentCodeInfo.AttempsInfo.Count > 5) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.ToManyAttemps, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.Code == req.CodeFromEmail) {
                var userId = await insertUser(_db, sentCodeInfo.User);
                _memoryCache.Remove(key);
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.ValidCode, Email: sentCodeInfo.User.Email, Name: sentCodeInfo.User.ClientName, ClientId: userId);
            } else {
                _memoryCache.Set(key, sentCodeInfo with {
                    AttempsInfo = sentCodeInfo.AttempsInfo with {
                        Count = sentCodeInfo.AttempsInfo.Count + 1
                    }
                }, new MemoryCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10)});

                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.IncorrectCode, Email: null, Name: null, ClientId: null);
            }
        }

        public async Task<RegistrateVerifyEmailCodeRes> VerifyEmailLink(RegistrateVerifyEmailLinkReq req){
            var tokenData = RegisterLinkToken.FromJson(_dataProtector.Unprotect(req.TokenFromRegisterLink));

            var key = $"reg_{tokenData.Email}";
            var sentCodeInfo = _memoryCache.Get<SentCodeInfo>(key);
            if (sentCodeInfo == null) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }
            if(req.DeviceId != sentCodeInfo.User.DeviceId) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.DeviceIdNotMatch, Email: null, Name: null, ClientId: null);
            }

            var codeExpired = sentCodeInfo.Until < DateTime.Now;
            if (codeExpired) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.NoSentCode, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.AttempsInfo.Count > 5) {
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.ToManyAttemps, Email: null, Name: null, ClientId: null);
            }

            if (sentCodeInfo.Code == tokenData.Code) {
                var userId = await insertUser(_db, sentCodeInfo.User);
                _memoryCache.Remove(key);
                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.ValidCode, Email: sentCodeInfo.User.Email, Name: sentCodeInfo.User.ClientName, ClientId: userId);
            } else {
                _memoryCache.Set(key, sentCodeInfo with {
                    AttempsInfo = sentCodeInfo.AttempsInfo with {
                        Count = sentCodeInfo.AttempsInfo.Count + 1
                    }
                }, new MemoryCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(10) });

                return new RegistrateVerifyEmailCodeRes(RegistrateVerifyEmailCodeResStatus.IncorrectCode, Email: null, Name: null, ClientId: null);
            }
        }


        public async Task<TryRemoveUserRes> TryRemoveUser(Guid userId) {
            var user = MobileUserModel.GetOrNull(userId, _db, _memoryCache);

            if (user == null) {
                return new TryRemoveUserRes(TryRemoveUserResStatus.UserNotExists);
            }

            await MobileUserModel.RemoveUser(userId, _db, _memoryCache);
            return new TryRemoveUserRes(TryRemoveUserResStatus.Ok);
        }

        private static async Task<Guid> insertUser(ZebraApplicationContext db, NeedToCreateUserInfo user) {
            var mobileUser = new MobileUserModel() {
                Name = user.ClientName,
                Email = user.Email,
                BirthDate = user.BirthDate.HasValue ? user.BirthDate.Value.ToUniversalTime() : null,
                RegDate = DateTime.UtcNow,
                Discount = 10,
                ZebraCoinBalance = 0,
                Status = MobileUserModel.UserStatus.Active
            };

            return await MobileUserModel.InsertUser(db, mobileUser);
        }
      

        public record SentCodeInfo(string Code, DateTime Until, NeedToCreateUserInfo User, VerifyEmailCodeAttempsInfo AttempsInfo);
        public record VerifyEmailCodeAttempsInfo(int Count);
        public record RegisterLinkToken(string Email, string Code) {
            public string ToJson() {
                return JsonConvert.SerializeObject(this);
            }
            public static RegisterLinkToken FromJson(string json) => JsonConvert.DeserializeObject<RegisterLinkToken>(json)!;
        };

        private static string randomNumbers(int length) {
            var random = new Random();
            var sb = new StringBuilder();
            for (var i = 0; i < length; i++) {
                sb.Append(random.Next(0, 9));
            }

            return sb.ToString();
        }
    }


    public interface IRegistrationEmailSenderService {
        public Task SendVerifyCode(string toEmail, string code, string authByLinkToken);
    }

    public class RegistrationEmailSenderService : IRegistrationEmailSenderService {
        private readonly EmailSenderConfig _conf;

        public RegistrationEmailSenderService(IOptions<EmailSenderConfig> conf) {
            _conf = conf.Value;
        }

        public async Task SendVerifyCode(string toEmail, string code, string authByLinkToken) {
            var smtpClient = new SmtpClient("smtp.mail.ru") {
                Port = 587,
                Credentials = new NetworkCredential(_conf.UserName, _conf.UserPwd),
                EnableSsl = true,
            };

            var continueOnMobileLink = $"{_conf.RegisterContinueOnMobileLink}/mobile-app?Type=registrate&amp;AuthByLinkToken={authByLinkToken}";

            var from = new MailAddress(_conf.UserName, "ZebraCoffee");
            var to = new MailAddress(toEmail);
            var m = new MailMessage(from, to);
            m.Subject = "Регистрация";
            m.Body = EmailMsgs.GetEmailVerifyMsg(code, continueOnMobileLink).Value;
            m.IsBodyHtml = true;

            smtpClient.Send(m);
        }
    }
}
