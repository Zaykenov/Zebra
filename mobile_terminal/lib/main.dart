import 'dart:async';
import 'dart:io';

import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_downloader/flutter_downloader.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_terminal/zebra_web_view.dart';

class MyHttpOverrides extends HttpOverrides {
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback =
          (X509Certificate cert, String host, int port) => true;
  }
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  if (!Platform.isWindows) {
    await FlutterDownloader.initialize(debug: true, ignoreSsl: false);
  }

  if (kDebugMode) {
    HttpOverrides.global = MyHttpOverrides();
  }

  // try {
  //   await Firebase.initializeApp();

  //   final fbm = FirebaseMessaging.instance;

  //   await fbm.requestPermission();
  //   await fbm.subscribeToTopic('CarGoRuqsatMobileTopic');

  //   await FbmHelper().checkCookiesAndSubscribeOrUnSubscribeToTopics();
  // } catch (e) {
  //   //var o = 0;
  // }

  runApp(MaterialApp(
      debugShowCheckedModeBanner: false,
      home: const StartScreen(),
      navigatorObservers: [FlutterSmartDialog.observer],
      builder: FlutterSmartDialog.init()));
}

class StartScreen extends StatefulWidget {
  const StartScreen({Key? key}) : super(key: key);

  @override
  StartScreenState createState() => StartScreenState();
}

class StartScreenState extends State<StartScreen> {
  //static const String initialUrl = "https://zebra-crm.kz/terminal/order";
  static const String initialUrl = "https://zebra.korsetu.kz/terminal/order";

  static const String internetConnectionStatusUrl =
      'https://zebra-api.korsetu.kz/ping';
  static const String createCheckUrl =
      'https://zebra-api.korsetu.kz/check/create';
  static const String saveFailedCheckUrl =
      'https://zebra-api.korsetu.kz/check/failed';

  @override
  void initState() {
    super.initState();

    WidgetsBinding.instance.addPostFrameCallback((_) async {
      try {
        // FirebaseMessaging.onBackgroundMessage((message) async {
        //   print(message.notification?.body);
        // });
        FirebaseMessaging.onMessage.listen((RemoteMessage message) {
          showDialog(
              context: context,
              builder: (BuildContext context) {
                return AlertDialog(
                  title: Text(message.notification?.title ?? "Нет заголовка"),
                  content: Text(message.notification?.body ?? "Нет содержания"),
                  actions: [
                    TextButton(
                      child: const Text("Ok"),
                      onPressed: () {
                        Navigator.of(context).pop();
                      },
                    )
                  ],
                );
              });
        });

        FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
          showDialog(
              context: context,
              builder: (BuildContext context) {
                return AlertDialog(
                  title: Text(message.notification?.title ?? "Нет заголовка"),
                  content: Text(message.notification?.body ?? "Нет содержания"),
                  actions: [
                    TextButton(
                      child: const Text("Ok"),
                      onPressed: () {
                        Navigator.of(context).pop();
                      },
                    )
                  ],
                );
              });
        });
      } catch (e) {
        //svar o = 0;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<bool>(
        future: isUserAuthorizedAndCanAuthWithBiometricts(),
        builder: (BuildContext context, AsyncSnapshot<bool> snapshot) {
          Widget ret;
          if (!snapshot.hasData) {
            ret = const Scaffold(
              body: Center(
                child: Text("Loading..."),
              ),
            );
          } else {
            //ret = const PrinterSettings();
            var isUserAuthorized = snapshot.requireData;
            if (isUserAuthorized) {
              // ret = const BiometryAuthScreen();
              ret = const Text("load..");
            } else {
              ret = const ZebraWebView(
                initialUrl: initialUrl,
                internetConnectionStatusUrl: internetConnectionStatusUrl,
                createCheckUrl: createCheckUrl,
                saveFailedCheckUrl: saveFailedCheckUrl,
              );
            }
          }

          return MaterialApp(
            debugShowCheckedModeBanner: false,
            home: ret,
          );
        });
  }

  Future<bool> isUserAuthorizedAndCanAuthWithBiometricts() async {
    return false;
  }
}
