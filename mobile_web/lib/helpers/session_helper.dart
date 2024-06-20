import 'dart:convert';

import 'package:shared_preferences/shared_preferences.dart';

import 'api_caller.dart';

class SessionHelper {
  Future<bool> isCurrentUserAuthenticated() async {
    var prefs = await SharedPreferences.getInstance();
    return prefs.getBool("isCurrentUserAuthenticated") ?? false;
  }

  Future userSignIn(String userEmail, String userName, String clientId) async {
    var prefs = await SharedPreferences.getInstance();
    await prefs.setBool("isCurrentUserAuthenticated", true);
    await prefs.setString("UserName", userName);
    await prefs.setString("UserEmail", userEmail);
    await prefs.setString("ClientId", clientId);
  }

  Future userSignOut() async {
    var prefs = await SharedPreferences.getInstance();
    await prefs.setBool("isCurrentUserAuthenticated", false);
    await prefs.remove("UserName");
    await prefs.remove("UserEmail");
    await prefs.remove("ClientId");
  }

  Future<String> getClientId() async {
    var prefs = await SharedPreferences.getInstance();
    return prefs.getString("ClientId")!;
  }

  Future<String> getClientUserName() async {
    var prefs = await SharedPreferences.getInstance();
    return prefs.getString("UserName")!;
  }

  Future<String?> getLastAuthLinkToken() async {
    var prefs = await SharedPreferences.getInstance();
    return prefs.getString("AuthLinkToken");
  }

  Future setLastAuthLinkToken(String authLinkToken) async {
    var prefs = await SharedPreferences.getInstance();
    return prefs.setString("AuthLinkToken", authLinkToken);
  }

  Future<AppConfig> getAppConfig() async {
    var prefs = await SharedPreferences.getInstance();

    var json = prefs.getString("AppConf");
    if (json != null) {
      return AppConfig.fromJson(jsonDecode(json));
    } else {
      var newConfig = AppConfig(
          apiUrl: ApiCaller().prodApiUrl, envName: EnvName.production);
      await setAppConfig(newConfig);
      return newConfig;
    }
  }

  Future setAppConfig(AppConfig conf) async {
    var prefs = await SharedPreferences.getInstance();
    await prefs.setString("AppConf", jsonEncode(conf.toJson()));
  }

  Future<FeedbackForOrder?> getFeedbackForOrder(int checkId) async {
    var prefs = await SharedPreferences.getInstance();

    var json = prefs.getString("FeedbackForOrder_$checkId");
    if (json == null) {
      return null;
    }

    return FeedbackForOrder.fromJson(jsonDecode(json));
  }

  Future setFeedbackForOrder(FeedbackForOrder feedback) async {
    var prefs = await SharedPreferences.getInstance();
    await prefs.setString(
        "FeedbackForOrder_${feedback.checkid}", jsonEncode(feedback.toJson()));
  }
}

class AppConfig {
  EnvName? envName;
  String? apiUrl;

  AppConfig({this.envName, this.apiUrl});

  AppConfig.fromJson(Map<String, dynamic> json) {
    envName = EnvName.values.firstWhere((e) => e.toString() == json['envName']);
    apiUrl = json['apiUrl'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['envName'] = envName.toString();
    data['apiUrl'] = apiUrl;
    return data;
  }
}

enum EnvName { staging, production }
