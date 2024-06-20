import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:latlong2/latlong.dart';
import 'package:mobile_web/helpers/session_helper.dart';

class ApiCaller {
  final String prodApiUrl = "https://zebra-crm.kz:13930";
  final String stagingApiUrl = "https://zebra-mobile-api.korsetu.kz";

  Future<TryCreateUserResStatus> tryCreateUser(TryCreateUserReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/registration/try-create-user"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    String status = jsonDecode(response.body)["Status"];
    return TryCreateUserResStatus.values
        .firstWhere((e) => e.toString() == 'TryCreateUserResStatus.$status');
  }

  Future<CreateIncognitoUserRes> createIncognitoUser(
      CreateIncognitoUserReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/registration/create-incognito-user"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    return CreateIncognitoUserRes.fromJson(jsonDecode(response.body));
  }

  Future<TrySignInResStatus> trySignIn(TrySignInModel req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response = await http.post(Uri.parse("$apiUrl/auth/try-sign-in"),
        headers: {
          "content-type": "application/json",
          "accept": "application/json",
        },
        body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    String status = jsonDecode(response.body)["Status"];
    return TrySignInResStatus.values
        .firstWhere((e) => e.toString() == 'TrySignInResStatus.$status');
  }

  Future<RegistrateVerifyEmailCodeRes> verifyRegEmailCode(
      VerifyEmailCodeReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/registration/verify-email-code"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    return RegistrateVerifyEmailCodeRes.fromJson(jsonDecode(response.body));
  }

  Future<SignInVerifyEmailCodeRes> verifySignInEmailCode(
      VerifyEmailCodeReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/auth/verify-email-code"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    return SignInVerifyEmailCodeRes.fromJson(jsonDecode(response.body));
  }

  Future<VerifyEmailLinkRes> verifySignInEmailLink(
      VerifyEmailLinkReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/auth/verify-email-link"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    return VerifyEmailLinkRes.fromJson(jsonDecode(response.body));
  }

  Future<VerifyEmailLinkRes> verifyRegistrateEmailLink(
      VerifyEmailLinkReq req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/registration/verify-email-link"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    return VerifyEmailLinkRes.fromJson(jsonDecode(response.body));
  }

  Future<TryGenerateQrRes> genQrCode(String clientId) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response = await http
        .get(Uri.parse("$apiUrl/user-qr/try-generate?userId=$clientId"));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    var json = jsonDecode(response.body);

    return TryGenerateQrRes.fromJson(json);
  }

  Future<TryGetUserInfoRes> getClientInfo(String clientId) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response = await http
        .get(Uri.parse("$apiUrl/user-info/try-get-info?userId=$clientId"));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    var json = jsonDecode(response.body);
    return TryGetUserInfoRes.fromJson(json);
  }

  Future<TryGetLastOrdersRes> getCurrentOrdersRes(String clientId) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response = await http.get(
        Uri.parse("$apiUrl/user-orders/try-get-last-orders?userId=$clientId"));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    var json = jsonDecode(response.body);
    return TryGetLastOrdersRes.fromJson(json);
  }

  Future<GetZebraLocationsRes> getZebraLocations() async {
    // final response =
    //     await http.get(Uri.parse("$apiUrl/get-zebra-locations"));

    // if (response.statusCode != 200) {
    //   throw Exception('Failed to load');
    // }

    // var json = jsonDecode(response.body);

    var json = GetZebraLocationsRes([
      ZebraShopLocation(
          name: "​Улица Мухтара Ауэзова, 17",
          location: LatLng(51.168849, 71.4234124)),
      ZebraShopLocation(
          name: "​Улица Иманова, 3/1в киоск",
          location: LatLng(51.163358, 71.431049))
    ]).toJson();

    return GetZebraLocationsRes.fromJson(json);
  }

  Future<List<Map<String, dynamic>>> fetchZebraLocations() async {
    try {
      final response = await http
          .get(Uri.parse('https://zebra-mobile-api.korsetu.kz/shop/get-all'));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final locations = data['Data'];
        return List<Map<String, dynamic>>.from(locations);
      } else {
        throw Exception('Failed to load shop locations');
      }
    } catch (error) {
      throw Exception('Failed to fetch shop locations: $error');
    }
  }

  Future<TrySetFeedbackResStatus> setFeedback(FeedbackForOrder req) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response =
        await http.post(Uri.parse("$apiUrl/user-feedback/try-set-feedback"),
            headers: {
              "content-type": "application/json",
              "accept": "application/json",
            },
            body: jsonEncode(req.toJson()));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    String status = jsonDecode(response.body)["Status"];
    return TrySetFeedbackResStatus.values
        .firstWhere((e) => e.toString() == 'TrySetFeedbackResStatus.$status');
  }

  Future<TryRemoveUserResStatus> tryRemoveUser(String clientId) async {
    var appConfig = await SessionHelper().getAppConfig();
    var apiUrl = appConfig.apiUrl;

    final response = await http.post(
        Uri.parse("$apiUrl/registration/try-remove-user?userId=$clientId"));

    if (response.statusCode != 200) {
      throw Exception('Failed to load');
    }

    String status = jsonDecode(response.body)["Status"];
    return TryRemoveUserResStatus.values
        .firstWhere((e) => e.toString() == 'TryRemoveUserResStatus.$status');
  }

  Future<List<ProductItem>> fetchProducts(String shopId) async {
    final response = await http.get(Uri.parse(
        'https://zebra-mobile-api.korsetu.kz/shop/get-tech-carts/$shopId'));

    if (response.statusCode == 200) {
      final parsedData = json.decode(response.body);
      final items = parsedData['Data']['Items'] as List<dynamic>;

      return items
          .map<ProductItem>((item) => ProductItem.fromJson(item))
          .toList();
    } else {
      throw Exception('Failed to load products');
    }
  }
}

class ProductItem {
  final int itemId;
  final String itemName;
  final double price;
  final String imageBase64;
  final String? category; // Modify as needed
  final List<Nabor>? nabors; // List of Nabors

  ProductItem({
    required this.itemId,
    required this.itemName,
    required this.price,
    required this.imageBase64,
    this.nabors,
    this.category,
  });

  factory ProductItem.fromJson(Map<String, dynamic> json) {
    final naborsData = json['Nabors'] as List<dynamic>?;
    final naborsList =
        naborsData?.map<Nabor>((nabor) => Nabor.fromJson(nabor)).toList();

    return ProductItem(
      itemId: json['ItemId'],
      category: json['Category'],
      itemName: json['ItemName'],
      price: json['Price'].toDouble(),
      imageBase64: json['ImageUrl'],
      nabors: naborsList,
    );
  }
}

class Nabor {
  final int naborId;
  final String name;
  final String description;
  final int min;
  final int max;
  final double price;
  final List<Modificator> modificators; // List of Modificators

  Nabor({
    required this.naborId,
    required this.name,
    required this.description,
    required this.min,
    required this.max,
    required this.price,
    required this.modificators,
  });

  factory Nabor.fromJson(Map<String, dynamic> json) {
    final modificatorsData = json['Modificators'] as List<dynamic>;
    final modificatorsList = modificatorsData
        .map<Modificator>((modificator) => Modificator.fromJson(modificator))
        .toList();

    return Nabor(
      naborId: json['NaborId'],
      name: json['Name'],
      description: json['Description'],
      min: json['Min'],
      max: json['Max'],
      price: json['Price'].toDouble(),
      modificators: modificatorsList,
    );
  }
}

class Modificator {
  final int ingredientId;
  final String ingredientName;
  final String image;
  int taps;
  bool isPressed;
  // Modify as needed

  Modificator({
    required this.ingredientId,
    required this.ingredientName,
    required this.image,
    this.taps = 0,
    this.isPressed = false,
  });

  factory Modificator.fromJson(Map<String, dynamic> json) {
    return Modificator(
      ingredientId: json['IngredientId'],
      ingredientName: json['IngredientName'],
      image: json['Image'],
    );
  }
}

class TryCreateUserReq {
  final String email;
  final String clientName;
  final DateTime? birthDate;
  final String deviceId;

  TryCreateUserReq(this.email, this.clientName, this.birthDate, this.deviceId);
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['email'] = email;
    data['clientName'] = clientName;
    data['birthDate'] = birthDate?.toIso8601String();
    data['deviceId'] = deviceId;
    return data;
  }
}

class CreateIncognitoUserReq {
  final String clientName;
  final String deviceId;

  CreateIncognitoUserReq(this.clientName, this.deviceId);
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['clientName'] = clientName;
    data['deviceId'] = deviceId;
    return data;
  }
}

class CreateIncognitoUserRes {
  late String? clientId;

  CreateIncognitoUserRes.fromJson(Map<String, dynamic> json) {
    clientId = json["ClientId"];
  }
}

enum TryCreateUserResStatus {
  // ignore: constant_identifier_names
  CodeSentToEmail,
  // ignore: constant_identifier_names
  AlreadySentToEmail,
  // ignore: constant_identifier_names
  AlreadyRegistered
}

class VerifyEmailCodeReq {
  final String email;
  final String codeFromEmail;
  final String deviceId;

  VerifyEmailCodeReq(this.email, this.codeFromEmail, this.deviceId);
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['email'] = email;
    data['codeFromEmail'] = codeFromEmail;
    data['deviceId'] = deviceId;
    return data;
  }
}

class VerifyEmailLinkReq {
  final String tokenFromRegisterLink;
  final String deviceId;

  VerifyEmailLinkReq(this.tokenFromRegisterLink, this.deviceId);
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['tokenFromRegisterLink'] = tokenFromRegisterLink;
    data['deviceId'] = deviceId;
    return data;
  }
}

class RegistrateVerifyEmailCodeRes {
  late RegistrateVerifyEmailCodeResStatus status;
  late String? clientId;

  RegistrateVerifyEmailCodeRes.fromJson(Map<String, dynamic> json) {
    status = RegistrateVerifyEmailCodeResStatus.values.firstWhere((e) =>
        e.toString() == 'RegistrateVerifyEmailCodeResStatus.${json["Status"]}');
    clientId = json["ClientId"];
  }
}

enum RegistrateVerifyEmailCodeResStatus {
  // ignore: constant_identifier_names
  ValidCode,
  // ignore: constant_identifier_names
  IncorrectCode,
  // ignore: constant_identifier_names
  ToManyAttemps,
  // ignore: constant_identifier_names
  NoSentCode,
  // ignore: constant_identifier_names
  DeviceIdNotMatch,
}

class SignInVerifyEmailCodeRes {
  late SignInVerifyEmailCodeResStatus status;
  late String? clientId;
  late String? name;
  late String? email;

  SignInVerifyEmailCodeRes.fromJson(Map<String, dynamic> json) {
    status = SignInVerifyEmailCodeResStatus.values.firstWhere((e) =>
        e.toString() == 'SignInVerifyEmailCodeResStatus.${json["Status"]}');
    clientId = json["ClientId"];
    name = json["Name"];
    email = json["Email"];
  }
}

enum SignInVerifyEmailCodeResStatus {
  // ignore: constant_identifier_names
  ValidCode,
  // ignore: constant_identifier_names
  IncorrectCode,
  // ignore: constant_identifier_names
  ToManyAttemps,
  // ignore: constant_identifier_names
  NoSentCode,
  // ignore: constant_identifier_names
  DeviceIdNotMatch,
}

class VerifyEmailLinkRes {
  late VerifyEmailLinkResStatus status;
  late String? email;
  late String? clientId;
  late String? name;

  VerifyEmailLinkRes.fromJson(Map<String, dynamic> json) {
    status = VerifyEmailLinkResStatus.values.firstWhere(
        (e) => e.toString() == 'VerifyEmailLinkResStatus.${json["Status"]}');
    email = json["Email"];
    clientId = json["ClientId"];
    name = json["Name"];
  }
}

enum VerifyEmailLinkResStatus {
  // ignore: constant_identifier_names
  ValidCode,
  // ignore: constant_identifier_names
  IncorrectCode,
  // ignore: constant_identifier_names
  ToManyAttemps,
  // ignore: constant_identifier_names
  NoSentCode,
  // ignore: constant_identifier_names
  DeviceIdNotMatch,
}

enum TrySignInResStatus {
  // ignore: constant_identifier_names
  CodeSentToEmail,
  // ignore: constant_identifier_names
  AlreadySentToEmail,
  // ignore: constant_identifier_names
  UserNotExists
}

class TrySignInModel {
  final String email;
  final String deviceId;

  TrySignInModel(this.email, this.deviceId);
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['email'] = email;
    data['deviceId'] = deviceId;
    return data;
  }
}

class TryGenerateQrRes {
  late TryGenerateQrResStatus status;
  late String? qrContent;
  late DateTime? until;

  TryGenerateQrRes.fromJson(Map<String, dynamic> json) {
    status = TryGenerateQrResStatus.values.firstWhere(
        (e) => e.toString() == 'TryGenerateQrResStatus.${json["Status"]}');
    qrContent = json["QrContent"];
    until = DateTime.tryParse(json["Until"]);
  }
}

enum TryGenerateQrResStatus {
  // ignore: constant_identifier_names
  UserNotExists,
  // ignore: constant_identifier_names
  Ok,
}

class TryGetUserInfoRes {
  late TryGetUserInfoResStatus status;
  late UserInfoRes? userInfo;

  TryGetUserInfoRes.fromJson(Map<String, dynamic> json) {
    status = TryGetUserInfoResStatus.values.firstWhere(
        (e) => e.toString() == 'TryGetUserInfoResStatus.${json["Status"]}');
    var userInfoJson = json["Info"];
    if (userInfoJson != null) {
      userInfo = UserInfoRes.fromJson(userInfoJson);
    } else {
      userInfo = null;
    }
  }
}

class UserInfoRes {
  late double? zebraCoinBalance;
  late double? discount;

  UserInfoRes.fromJson(Map<String, dynamic> json) {
    zebraCoinBalance = json["ZebraCoinBalance"];
    discount = json["Discount"] * 100;
  }
}

enum TryGetUserInfoResStatus {
  // ignore: constant_identifier_names
  UserNotExists,
  // ignore: constant_identifier_names
  Ok,
}

class GenQrAuthCodeRes {
  final String authCode;
  final DateTime until;

  GenQrAuthCodeRes(this.authCode, this.until);
}

class GetZebraLocationsRes {
  List<ZebraShopLocation>? locations;

  GetZebraLocationsRes(this.locations);

  GetZebraLocationsRes.fromJson(Map<String, dynamic> json) {
    if (json['locations'] != null) {
      locations = <ZebraShopLocation>[];
      json['locations'].forEach((v) {
        locations!.add(ZebraShopLocation.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['locations'] =
        locations != null ? locations!.map((v) => v.toJson()).toList() : null;
    return data;
  }
}

class ZebraShopLocation {
  LatLng? location;
  String? name;

  ZebraShopLocation({this.location, this.name});

  ZebraShopLocation.fromJson(Map<String, dynamic> json) {
    name = json['Name'];
    location = LatLng(json['Latitude'], json['Longitude']);
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Name'] = name;
    data['Latitude'] = location?.latitude;
    data['Longitude'] = location?.longitude;
    return data;
  }
}

class TryGetLastOrdersRes {
  late TryGetLastOrdersResStatus status;
  late List<Order>? orders;

  TryGetLastOrdersRes.fromJson(Map<String, dynamic> json) {
    status = TryGetLastOrdersResStatus.values.firstWhere(
        (e) => e.toString() == 'TryGetLastOrdersResStatus.${json["Status"]}');

    orders = <Order>[];
    if (json["Orders"] != null) {
      json['Orders'].forEach((v) {
        orders!.add(Order.fromJson(v));
      });
    }
  }
}

enum TryGetLastOrdersResStatus {
  // ignore: constant_identifier_names
  UserNotExists,
  // ignore: constant_identifier_names
  Ok,
}

class Order {
  late int id;
  late DateTime openedat;
  late double sum;
  late String payment;
  late List<OrderTovar> tovars;
  late List<OrderTechCart> techCarts;
  String? comment;
  String? feedback;

  Order.fromJson(Map<String, dynamic> json) {
    id = json['Id'];
    openedat = DateTime.parse(json['OpenedAt']);
    sum = double.parse(json['Sum'].toString());
    payment = json['Payment'];

    tovars = <OrderTovar>[];
    if (json['Tovars'] != null) {
      json['Tovars'].forEach((v) {
        tovars.add(OrderTovar.fromJson(v));
      });
    }
    techCarts = <OrderTechCart>[];
    if (json['TechCarts'] != null) {
      json['TechCarts'].forEach((v) {
        techCarts.add(OrderTechCart.fromJson(v));
      });
    }
    comment = json['Comment'];
    feedback = json['Feedback'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Id'] = id;
    data['OpenedAt'] = openedat.toIso8601String();
    data['Sum'] = sum;
    data['Payment'] = payment;
    data['Tovars'] = tovars.map((v) => v.toJson()).toList();
    data['TechCarts'] = techCarts.map((v) => v.toJson()).toList();
    data['Comment'] = comment;
    data['Feedback'] = feedback;
    return data;
  }
}

class OrderTovar {
  late String name;
  late int quantity;
  late double price;

  OrderTovar.fromJson(Map<String, dynamic> json) {
    name = json['Name'];
    quantity = json['Quantity'];
    price = json['Price'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Name'] = name;
    data['Quantity'] = quantity;
    data['Price'] = price;
    return data;
  }
}

class OrderTechCart {
  late String name;
  late int quantity;
  late double price;
  late List<String> modificators;

  OrderTechCart.fromJson(Map<String, dynamic> json) {
    name = json['Name'];
    quantity = json['Quantity'];
    price = json['Price'];

    modificators = <String>[];

    if (json['Modificators'] != null) {
      json['Modificators'].forEach((v) {
        modificators.add(v.toString());
      });
    }
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Name'] = name;
    data['Quantity'] = quantity;
    data['Price'] = price;
    data['Modificators'] = modificators;
    return data;
  }
}

class FeedbackForOrder {
  late int checkid;
  late String userId;
  late double scoreQuality;
  late double scoreService;
  late String feedbackText;
  late String checkJson;

  FeedbackForOrder(this.checkid, this.userId, this.scoreQuality,
      this.scoreService, this.feedbackText, this.checkJson);

  FeedbackForOrder.fromJson(Map<String, dynamic> json) {
    checkid = json['Checkid'];
    userId = json['UserId'];
    scoreQuality = json['ScoreQuality'];
    scoreService = json['ScoreService'];
    feedbackText = json['FeedbackText'];
    checkJson = json['CheckJson'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['Checkid'] = checkid;
    data['UserId'] = userId;
    data['ScoreQuality'] = scoreQuality;
    data['ScoreService'] = scoreService;
    data['FeedbackText'] = feedbackText;
    data['CheckJson'] = checkJson;
    return data;
  }
}

enum TrySetFeedbackResStatus {
  // ignore: constant_identifier_names
  UserNotExists,
  // ignore: constant_identifier_names
  Ok
}

enum TryRemoveUserResStatus {
  // ignore: constant_identifier_names
  UserNotExists,
  // ignore: constant_identifier_names
  Ok
}

class OrderAndFeedback {
  Order order;
  FeedbackForOrder? feedback;

  OrderAndFeedback(this.order, this.feedback);
}
