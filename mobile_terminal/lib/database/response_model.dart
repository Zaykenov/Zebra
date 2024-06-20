import 'dart:convert';

class Response {
  int id;
  String time;
  String request;
  String response;
  String exception;
  int retryCount;
  String status;

  Response({
    required this.id,
    required this.time,
    required this.request,
    required this.response,
    required this.exception,
    required this.retryCount,
    required this.status,
  });

  factory Response.fromJson(Map<String, dynamic> json) => Response(
        id: json["id"],
        time: json["time"],
        request: json["request"],
        response: json["response"],
        exception: json["exception"],
        retryCount: json["retryCount"],
        status: json["status"],
      );

  Map<String, dynamic> toJson() => {
        "id": id,
        "time": time,
        "request": request,
        "response": response,
        "exception": exception,
        "retryCount": retryCount,
        "status": status,
      };

  static Response clientFromJson(String str) {
    final jsonData = json.decode(str);
    return Response.fromJson(jsonData);
  }

  static String clientToJson(Response data) {
    final dyn = data.toJson();
    return json.encode(dyn);
  }
}
