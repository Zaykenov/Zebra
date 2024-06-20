namespace ClickHouseProxy {
    public record VisitInfo(
        string Type,
        string VisitId,
        string Brouser,
        DateTime Date,
        string ForwardedFor,
        string RemoteAddr,
        string Referer,
        string Url,
        string Method,
        string Body,
        string StatusCode
    );
}
