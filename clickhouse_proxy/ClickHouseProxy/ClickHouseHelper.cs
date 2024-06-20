using Microsoft.Extensions.Primitives;

namespace ClickHouseProxy {
    public static class ClickHouseHelper {
        public static string CreateTable = "CREATE TABLE IF NOT EXISTS clickhouse_logs(type String, id String, brouser String, timestamp DateTime, forwardedFor String, remoteAddr String, referer String, url String, method String, body String, statusCode String) engine Memory";

        public static string GetForwardedIpOrNull(HttpContext context) {
            var ip = GetHeaderValueAs<string>("X-Forwarded-For", context).SplitCsv().FirstOrDefault();
            if (ip.IsNullOrWhitespace()) {
                return null; // :(
            }
            return ip;
        }
        public static string GetRawUrl(this HttpRequest request) {
            return Microsoft.AspNetCore.Http.Extensions.UriHelper.GetDisplayUrl(request);
        }
        public static T GetHeaderValueAs<T>(string headerName, HttpContext context) {
            StringValues values;

            if (context?.Request?.Headers?.TryGetValue(headerName, out values) ?? false) {
                string rawValues = values.ToString();   // writes out as Csv when there are multiple.

                if (!rawValues.IsNullOrWhitespace())
                    return (T)Convert.ChangeType(values.ToString(), typeof(T));
            }
            return default(T);
        }
        public static List<string> SplitCsv(this string csvList, bool nullOrWhitespaceInputReturnsNull = false) {
            if (string.IsNullOrWhiteSpace(csvList))
                return nullOrWhitespaceInputReturnsNull ? null : new List<string>();

            return csvList
                .TrimEnd(',')
                .Split(',')
                .AsEnumerable<string>()
                .Select(s => s.Trim())
                .ToList();
        }
        public static bool IsNullOrWhitespace(this string s) {
            return string.IsNullOrWhiteSpace(s);
        }
    }
}
