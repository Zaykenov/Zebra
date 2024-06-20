import 'dart:async';

import 'package:app_links/app_links.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/app_link_helper.dart';
import 'package:mobile_web/main.dart';
import 'package:mobile_web/qr_page.dart';
import 'package:pinput/pinput.dart';
import 'package:platform_device_id/platform_device_id.dart';

import 'helpers/session_helper.dart';

class RegistrationValidationPage extends StatefulWidget {
  final String email;
  final String userName;
  const RegistrationValidationPage(this.email, this.userName, {Key? key})
      : super(key: key);

  @override
  State<RegistrationValidationPage> createState() =>
      _RegistrationValidationPageState();
}

class _RegistrationValidationPageState
    extends State<RegistrationValidationPage> {
  late AppLinks _appLinks;
  StreamSubscription<Uri>? _linkSubscription;

  @override
  void initState() {
    super.initState();

    initDeepLinks();
  }

  Future<void> initDeepLinks() async {
    _appLinks = AppLinks();

    // Check initial link if app was in cold state (terminated)
    final appLink = await _appLinks.getInitialAppLink();
    if (appLink != null) {
      if (kDebugMode) {
        print('getInitialAppLink: $appLink');
      }
      var res = await AppLinkHelper().proccessAppDeepLink(appLink);
      if (res == NeedToRestartStatus.NeedRestart) {
        if (mounted) {
          RestartWidget.restartApp(context);
        }
      }
    }

    // Handle link when app is in warm state (front or background)
    _linkSubscription = _appLinks.uriLinkStream.listen((uri) async {
      if (kDebugMode) {
        print('onAppLink: $uri');
      }
      var res = await AppLinkHelper().proccessAppDeepLink(uri);
      if (res == NeedToRestartStatus.NeedRestart) {
        if (mounted) {
          RestartWidget.restartApp(context);
        }
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: SingleChildScrollView(
      padding: const EdgeInsets.symmetric(vertical: 100),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          const Image(image: AssetImage('assets/logo.png')),
          Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              children: [
                Card(
                  color: const Color(0xff3EB2B2),
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Text(
                      "На ${widget.email} был отправлен код подтверждения",
                      textAlign: TextAlign.center,
                      style: const TextStyle(
                          color: Colors.white, fontStyle: FontStyle.italic),
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 50),
            child: Column(
              children: [
                Pinput(
                  autofocus: true,
                  defaultPinTheme: defaultPinTheme,
                  onCompleted: (emailCode) async {
                    String? deviceId = await PlatformDeviceId.getDeviceId;

                    SmartDialog.showLoading(msg: "Осталось совсем чуть-чуть");
                    var res = await ApiCaller().verifyRegEmailCode(
                        VerifyEmailCodeReq(widget.email, emailCode, deviceId!));
                    await SmartDialog.dismiss();

                    if (res.status ==
                        RegistrateVerifyEmailCodeResStatus.IncorrectCode) {
                      await SmartDialog.showNotify(
                          msg: "Некорректный код",
                          notifyType: NotifyType.error);
                    }
                    if (res.status ==
                        RegistrateVerifyEmailCodeResStatus.ToManyAttemps) {
                      SmartDialog.showNotify(
                          msg: "Слишком много попыток. Попробуйте позже",
                          notifyType: NotifyType.error);
                    }
                    if (res.status ==
                        RegistrateVerifyEmailCodeResStatus.ValidCode) {
                      await SessionHelper().userSignIn(
                          widget.email, widget.userName, res.clientId!);

                      SmartDialog.showToast("Вы успешно зарегистрировались!");
                      if (mounted) {
                        await Navigator.pushAndRemoveUntil(
                            context,
                            MaterialPageRoute(
                                builder: (context) => const QrPage()),
                            (route) => false);
                      }
                    }
                  },
                ),
              ],
            ),
          ),
        ],
      ),
    ));
  }

  final defaultPinTheme = PinTheme(
    width: 60,
    height: 60,
    textStyle: const TextStyle(
        fontSize: 20, color: Colors.black, fontWeight: FontWeight.w600),
    decoration: BoxDecoration(
        border: Border.all(color: const Color.fromRGBO(234, 239, 243, 1)),
        borderRadius: BorderRadius.circular(20),
        color: const Color.fromRGBO(234, 239, 243, 1)),
  );
}
