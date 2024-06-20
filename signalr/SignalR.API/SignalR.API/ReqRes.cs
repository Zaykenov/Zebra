using SignalR.API.Hubs;

namespace SignalR.API {
    public record PostNfcDiscountReq(string NfcContent, string UserMobileCode);
    public record TryPostNfcDiscountRes(TryPostNfcDiscountResStatus Status, string Message);

    [Newtonsoft.Json.JsonConverter(typeof(Newtonsoft.Json.Converters.StringEnumConverter))]
    [System.Text.Json.Serialization.JsonConverter(typeof(System.Text.Json.Serialization.JsonStringEnumConverter))]
    public enum TryPostNfcDiscountResStatus {
        Fail = 1,
        Ok = 2,
    }
}
