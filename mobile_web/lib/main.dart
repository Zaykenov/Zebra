import 'dart:async';
import 'dart:io';
import 'dart:ui';

import 'package:app_links/app_links.dart';
import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/app_link_helper.dart';
import 'package:mobile_web/qr_page.dart';
import 'package:mobile_web/remove_app_page.dart';
import 'package:mobile_web/shoping/cart/cart_item.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:shopping_cart/shopping_cart.dart';

import 'choose_signin_type_page.dart';
import 'helpers/session_helper.dart';

import 'package:flutter_background_service/flutter_background_service.dart';
import 'package:flutter_background_service_android/flutter_background_service_android.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

// class MyHttpOverrides extends HttpOverrides {
//   @override
//   HttpClient createHttpClient(SecurityContext? context) {
//     return super.createHttpClient(context)
//       ..badCertificateCallback =
//           (X509Certificate cert, String host, int port) => true;
//   }
// }
//

/*
void main() {
  scheduleTimeout(5 * 1000); // 5 seconds.
}

Timer scheduleTimeout([int milliseconds = 10000]) =>
    Timer(Duration(milliseconds: milliseconds), handleTimeout);

void handleTimeout() {  // callback function
  // Do some work.
}
*/
Future<void> main() async {
  //HttpOverrides.global = MyHttpOverrides();

  // ShoppingCart.init<CartItem>();
  // await initializeService();

  // runApp(const RestartWidget(child: ZebraApp()));
  runApp(const RestartWidget(child: RemoveAppPage()));
}

class ZebraApp extends StatefulWidget {
  const ZebraApp({super.key});

  @override
  State<ZebraApp> createState() => _ZebraAppState();
}

class _ZebraAppState extends State<ZebraApp> {
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
    return FutureBuilder(
      builder: (ctx, isCurrentUserAuthenticated) {
        Widget ret;

        if (!isCurrentUserAuthenticated.hasData) {
          ret = const Scaffold(
            body: Center(
                child: Image(
              image: AssetImage("assets/logo.png"),
            )),
          );
        } else {
          ret = isCurrentUserAuthenticated.requireData
              ? const QrPage()
              : const ChooseSignInTypePage();
        }

        return MaterialApp(
          debugShowCheckedModeBanner: false,
          theme: ThemeData(
              //primarySwatch: Colors.blue,
              primaryColor: const Color(0xff3EB2B2),
              colorScheme: const ColorScheme.light(primary: Color(0xff3EB2B2)),
              fontFamily: 'Kameron',
              useMaterial3: true),
          home: ret,
          navigatorObservers: [FlutterSmartDialog.observer],
          builder: FlutterSmartDialog.init(),
          //locale: Locale("ru", "RU"),
        );
      },
      future: SessionHelper().isCurrentUserAuthenticated(),
    );
  }
}

class RestartWidget extends StatefulWidget {
  const RestartWidget({super.key, this.child});

  final Widget? child;

  static void restartApp(BuildContext context) {
    context.findAncestorStateOfType<_RestartWidgetState>()?.restartApp();
  }

  @override
  State<StatefulWidget> createState() {
    return _RestartWidgetState();
  }
}

class _RestartWidgetState extends State<RestartWidget> {
  Key key = UniqueKey();

  void restartApp() {
    setState(() {
      key = UniqueKey();
    });
  }

  @override
  Widget build(BuildContext context) {
    return KeyedSubtree(
      key: key,
      child: widget.child ?? Container(),
    );
  }
}
