import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:platform_device_id/platform_device_id.dart';

class AppLinkHelper {
  Future<NeedToRestartStatus> proccessAppDeepLink(Uri uri) async {
    var type = uri.queryParameters["Type"];
    if (type == "sign-in") {
      return proccessSignInDeepLink(uri);
    } else {
      if (type == "registrate") {
        return proccessRegDeepLink(uri);
      }
    }

    return NeedToRestartStatus.NoNeedToRestart;
  }

  Future<NeedToRestartStatus> proccessSignInDeepLink(Uri uri) async {
    var authByLinkToken = uri.queryParameters["AuthByLinkToken"];
    var lastToken = await SessionHelper().getLastAuthLinkToken();
    if (lastToken == authByLinkToken) {
      // SmartDialog.showToast(
      //     "повторно ссылку использовать нельзя(\nпопробуйте ввести код");
      return NeedToRestartStatus.NoNeedToRestart;
    }

    SmartDialog.showLoading(msg: "Осталось совсем чуть-чуть");
    String? deviceId = await PlatformDeviceId.getDeviceId;
    var res = await ApiCaller()
        .verifySignInEmailLink(VerifyEmailLinkReq(authByLinkToken!, deviceId!));
    await SmartDialog.dismiss();

    if (res.status == VerifyEmailLinkResStatus.ValidCode) {
      await SessionHelper().setLastAuthLinkToken(authByLinkToken);
      await SessionHelper().userSignIn(res.email!, res.name!, res.clientId!);
      SmartDialog.showToast("Добро пожаловать!");
      return NeedToRestartStatus.NeedRestart;
    } else {
      SmartDialog.showToast("Ссылка перестала быть действительной(");
      return NeedToRestartStatus.NoNeedToRestart;
    }
  }

  Future<NeedToRestartStatus> proccessRegDeepLink(Uri uri) async {
    var authByLinkToken = uri.queryParameters["AuthByLinkToken"];
    var lastToken = await SessionHelper().getLastAuthLinkToken();
    if (lastToken == authByLinkToken) {
      // SmartDialog.showToast(
      //     "повторно ссылку использовать нельзя(\nпопробуйте ввести код");
      return NeedToRestartStatus.NoNeedToRestart;
    }

    SmartDialog.showLoading(msg: "Осталось совсем чуть-чуть");
    String? deviceId = await PlatformDeviceId.getDeviceId;
    var res = await ApiCaller().verifyRegistrateEmailLink(
        VerifyEmailLinkReq(authByLinkToken!, deviceId!));
    await SmartDialog.dismiss();

    if (res.status == VerifyEmailLinkResStatus.ValidCode) {
      await SessionHelper().setLastAuthLinkToken(authByLinkToken);
      await SessionHelper().userSignIn(res.email!, res.name!, res.clientId!);
      SmartDialog.showToast("Добро пожаловать!");
      return NeedToRestartStatus.NeedRestart;
    } else {
      SmartDialog.showToast("Ссылка перестала быть действительной(");
      return NeedToRestartStatus.NoNeedToRestart;
    }
  }
}

enum NeedToRestartStatus {
  // ignore: constant_identifier_names
  NoNeedToRestart,
  // ignore: constant_identifier_names
  NeedRestart
}
