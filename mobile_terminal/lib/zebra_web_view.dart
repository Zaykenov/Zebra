import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_downloader/flutter_downloader.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:intl/intl.dart';
import 'package:mobile_terminal/cache/shared_pref_helper.dart';
import 'package:mobile_terminal/check_models/check_model.dart';
import 'package:mobile_terminal/database/database.dart';
import 'package:mobile_terminal/print_helper.dart';
import 'package:mobile_terminal/printer_settings.dart';
import 'package:path_provider/path_provider.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:http/http.dart' as http;
import 'package:cron/cron.dart';

class ZebraWebView extends StatefulWidget {
  final String initialUrl;
  final String internetConnectionStatusUrl;
  final String createCheckUrl;
  final String saveFailedCheckUrl;

  const ZebraWebView({
    Key? key,
    required this.initialUrl,
    required this.internetConnectionStatusUrl,
    required this.createCheckUrl,
    required this.saveFailedCheckUrl,
  }) : super(key: key);

  @override
  ZebraWebViewState createState() => ZebraWebViewState();
}

class ZebraWebViewState extends State<ZebraWebView> {
  InAppWebViewController? webView;
  late PullToRefreshController _c;
  String url = "";
  double progress = 0;

  StreamSubscription? _updateCoinsSubs;
  @override
  void initState() {
    super.initState();

    _c = PullToRefreshController(
        onRefresh: () async {
          if (Platform.isAndroid) {
            webView?.reload();
          } else if (Platform.isIOS) {
            webView?.loadUrl(
                urlRequest: URLRequest(url: await webView?.getUrl()));
          }

          await _c.endRefreshing();
        },
        options: PullToRefreshOptions(color: const Color(0xff3eb2b2)));
    Permission.storage.request();
    if (!Platform.isWindows) {
      FlutterDownloader.registerCallback(downloadCallback);
    }

    // startSendOrders();
    // startSendFailesOrders();

    //startCoinsTimer();
  }

  void startCoinsTimer() async {
    await cancelCoinsTimer();
    _updateCoinsSubs =
        Stream.periodic(const Duration(seconds: 5)).listen((event) async {
      if (mounted) {
        try {
          final responseNet = await http.get(
            //Uri.parse("https://www.google.kz/"),
            Uri.parse(widget.internetConnectionStatusUrl),
            //Uri.parse("https://zebra-crm.kz/login"),
          );

          if (responseNet.statusCode != 200) {
            return;
          }
        } catch (e) {
          var rr = e;
        }
      }
    });
  }

  Future<void> cancelCoinsTimer() async {
    await _updateCoinsSubs?.cancel();
  }

  Future<void> startSendOrders() async {
    final cron = Cron();
    cron.schedule(Schedule.parse('*/5 * * * * *'), () async {
      final responseNet = await http.get(
        Uri.parse(widget.internetConnectionStatusUrl),
      );

      if (responseNet.statusCode != 200) {
        return;
      }

      final checks = await DbHelper.db.getPendingJobs();
      for (var check in checks) {
        final String requestBody = check.request;
        final response = await http.post(
          Uri.parse(widget.createCheckUrl),
          headers: <String, String>{
            'Content-Type': 'application/json',
            'Authorization': check.authorization,
            'Idempotency-Key': check.idempotencyKey,
          },
          body: requestBody,
        );

        String status;
        if (response.statusCode == 200) {
          status = "Finished";
        } else {
          status = "Pending";
        }

        await DbHelper.db.updateJob(check.id, check.retryCount + 1, status,
            response.statusCode.toString());
      }
    });
  }

  Future<void> startSendFailesOrders() async {
    final cron = Cron();
    cron.schedule(Schedule.parse('*/20 * * * * *'), () async {
      final responseNet = await http.get(
        Uri.parse(widget.internetConnectionStatusUrl),
      );
      if (responseNet.statusCode != 200) {
        return;
      }

      final checks = await DbHelper.db.getFailedJobs();
      for (var check in checks) {
        final response = await http.post(
          Uri.parse(widget.saveFailedCheckUrl),
          headers: <String, String>{
            'Content-Type': 'application/json',
            'Authorization': check.authorization,
            'Idempotency-Key': check.idempotencyKey,
          },
          body: jsonEncode({
            'request': check.request, // Send request string
            'response': check.response // Send response string
          }),
        );

        if (response.statusCode == 200) {
          await DbHelper.db.deleteJob(check.id);
        } else {
          await DbHelper.db.updateJob(check.id, check.retryCount + 1, "Pending",
              response.statusCode.toString());
        }
      }
    });
  }

  @pragma('vm:entry-point')
  static void downloadCallback(
      String id, DownloadTaskStatus status, int progress) {}

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () => exitApp(),
      child: Scaffold(
        body: Column(children: <Widget>[
          Container(
            padding: const EdgeInsets.all(15),
            color: Colors.white,
          ),
          Container(
              padding: const EdgeInsets.all(15.0),
              color: Colors.white,
              child: progress < 1.0
                  ? LinearProgressIndicator(
                      value: progress,
                      backgroundColor: Colors.white,
                      valueColor: const AlwaysStoppedAnimation<Color>(
                          Color(0xff3eb2b2)),
                    )
                  : Container()),
          Expanded(
            child: InAppWebView(
              pullToRefreshController: _c,
              initialUrlRequest: URLRequest(
                  url: Uri.parse(widget.initialUrl),
                  headers: {'sender': 'mobile-app'}),
              initialOptions: InAppWebViewGroupOptions(
                  android: AndroidInAppWebViewOptions(
                      useHybridComposition: true,
                      useShouldInterceptRequest: false),
                  crossPlatform: InAppWebViewOptions(
                    //debuggingEnabled: true,
                    // useShouldInterceptFetchRequest: true,
                    supportZoom: false,
                    javaScriptEnabled: true,
                    cacheEnabled: true,
                    useShouldOverrideUrlLoading: true,
                    useOnDownloadStart: true,
                    incognito: false,
                    clearCache: false,
                    useOnLoadResource: true,
                    useShouldInterceptAjaxRequest: true,
                  ),
                  ios: IOSInAppWebViewOptions(
                    allowsBackForwardNavigationGestures: true,
                    sharedCookiesEnabled: false,
                  )),
              onWebViewCreated: (InAppWebViewController controller) {
                webView = controller;
              },
              onLoadResource: ((controller, resource) async {}),
              shouldInterceptAjaxRequest: (controller, ajaxRequest) async {
                var urlStr = ajaxRequest.url.toString();
                if (urlStr.endsWith("/terminal/order/print-receipt")) {
                  var printerOrNull =
                      await SharedPrefHelper().getSavedPrinterOrNull();
                  if (printerOrNull == null) {
                    SmartDialog.show(builder: (context) {
                      return Container(
                        height: 80,
                        width: 180,
                        decoration: BoxDecoration(
                          color: Colors.black,
                          borderRadius: BorderRadius.circular(10),
                        ),
                        alignment: Alignment.center,
                        child: const Text("Не настроен принтер!",
                            style: TextStyle(color: Colors.white)),
                      );
                    });
                  } else {
                    try {
                      var json =
                          jsonDecode(ajaxRequest.data)["data"].toString();

                      var model = CheckModel.fromJson(jsonDecode(json));

                      if (printerOrNull.usbPrinter != null) {
                        await PrintHelper().printUsbReceipt(
                            printerOrNull.usbPrinter!, model,
                            count: 2);
                      } else {
                        if (printerOrNull.wiFiPrinter != null) {
                          await PrintHelper().printWiFiReceipt(
                              printerOrNull.wiFiPrinter!, model);
                        }
                      }
                    } catch (e) {
                      SmartDialog.dismiss();
                      SmartDialog.show(builder: (context) {
                        return Container(
                          height: 80,
                          width: 180,
                          decoration: BoxDecoration(
                            color: Colors.black,
                            borderRadius: BorderRadius.circular(10),
                          ),
                          alignment: Alignment.center,
                          child: Text("Ошибка при распечатке: $e",
                              style: const TextStyle(color: Colors.white)),
                        );
                      });
                    }
                  }

                  ajaxRequest.url = Uri.parse(
                      "https://zebra-api.korsetu.kz/check/tag/getAll");
                } else if (urlStr.endsWith('/check/create')) {
                  final headers = ajaxRequest.headers!.getHeaders();

                  var job = SendOrderToApiJob(
                      0,
                      DateTime.now(),
                      ajaxRequest.data.toString(),
                      "",
                      "",
                      0,
                      "Pending",
                      headers['Authorization'],
                      headers['Idempotency-Key']);

                  await DbHelper.db.insertJob(job);
                }

                return ajaxRequest;
              },
              onLoadError: (controller, url, code, message) {
                //var i = 0;
              },
              onConsoleMessage: (controller, consoleMessage) {
                if (kDebugMode) {
                  print(consoleMessage.message);
                }
              },
              shouldOverrideUrlLoading: (controller, navigationAction) async {
                //var url = request.url;
                var uri = navigationAction.request.url!;
                var uriStr = uri.toString();

                if (uriStr.contains("t.me") ||
                    uriStr.contains("api.whatsapp.com") ||
                    uriStr.contains("mailto:")) {
                  await launch(uriStr);
                  return NavigationActionPolicy.CANCEL;
                }

                if (uri.host != 'zebra-crm.kz' &&
                    !uri.host.contains('korsetu.kz')) {
                  await launch(uriStr);
                  return NavigationActionPolicy.CANCEL;
                }

                if (![
                  "http",
                  "https",
                  "file",
                  "chrome",
                  "data",
                  "javascript",
                  "about"
                ].contains(uri.scheme)) {
                  if (await canLaunch(uriStr)) {
                    await launch(uriStr);
                    return NavigationActionPolicy.CANCEL;
                  }
                }

                // navigationAction.request.headers
                //     ?.addAll({'sender': 'mobile-app'});
                // return NavigationActionPolicy.ALLOW;

                navigationAction.request.headers ??= <String, String>{};
                if (navigationAction.request.headers!.containsKey("sender")) {
                  return NavigationActionPolicy.ALLOW;
                } else {
                  navigationAction.request.headers
                      ?.addAll({'sender': 'mobile-app'});
                  await controller.loadUrl(
                      urlRequest: navigationAction.request);
                  return NavigationActionPolicy.CANCEL;
                }
              },
              onLoadStop: (InAppWebViewController controller, Uri? url) async {
                setState(() {
                  this.url = url.toString();
                });
              },
              onProgressChanged:
                  (InAppWebViewController controller, int progress) {
                setState(() {
                  this.progress = progress / 100;
                });
              },
              onReceivedServerTrustAuthRequest:
                  (InAppWebViewController controller,
                      URLAuthenticationChallenge challenge) async {
                return ServerTrustAuthResponse(
                    action: ServerTrustAuthResponseAction.PROCEED);
              },
              onDownloadStart: (
                controller,
                url,
              ) async {
                await FlutterDownloader.enqueue(
                  url: url.toString(),
                  //fileName: '${DateTime.now().millisecondsSinceEpoch}',
                  savedDir: await getDownloadPath(),
                  showNotification: true,
                  openFileFromNotification: true,
                  saveInPublicStorage: false,
                );
              },
            ),
          ),
        ]),
        floatingActionButton: InkWell(
          onLongPress: () async {
            // if (kDebugMode) {
            //   var keysCount =
            //       await (await DbHelper.getInstance()).getCacheKeysCount();
            //   SmartDialog.showToast("keys count = $keysCount");
            // }
            SmartDialog.show(
              alignment: Alignment.center,
              //maskColor: Colors.transparent,
              builder: (_) {
                return Center(
                  child: Container(
                    height: 150,
                    width: 300,
                    color: Colors.white,
                    alignment: Alignment.center,
                    child: Card(
                      elevation: 10,
                      color: Colors.amber,
                      child: SizedBox(
                        height: 50,
                        child: TextButton(
                            onPressed: () async {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                    builder: (context) =>
                                        const PrinterSettings()),
                              );

                              SmartDialog.dismiss();
                            },
                            style: TextButton.styleFrom(
                                backgroundColor: Colors.blue),
                            child: const Text(
                              'Настройки принтера',
                              style: TextStyle(color: Colors.white),
                            )),
                      ),
                    ),
                  ),
                );
              },
            );
          },
          child: FloatingActionButton(
            // isExtended: true,
            backgroundColor: const Color(0xff3eb2b2),
            onPressed: () async {
              _reload();
            },
            child: const Icon(Icons.refresh),
          ),
        ),
      ),
    );
  }

  Future<void> _reload() async {
    if (Platform.isAndroid) {
      webView?.reload();
    } else if (Platform.isIOS) {
      webView?.loadUrl(urlRequest: URLRequest(url: await webView?.getUrl()));
    }
  }

  Future<String> getDownloadPath() async {
    Directory? directory;
    if (Platform.isIOS) {
      directory = await getApplicationDocumentsDirectory();
    } else {
      directory = await getExternalStorageDirectory();
      directory ??=
          await Directory('/storage/emulated/0/Download/zebra-crm').create();

      directory = await Directory(
              '${directory.path}/${DateFormat('dd-MM-yyyy HH:mm:ss').format(DateTime.now())}')
          .create();
    }
    return directory.path;
  }

  Future<bool> exitApp() async {
    if (webView == null) {
      return false;
    }

    if (await webView!.canGoBack()) {
      await webView!.goBack();
      return false;
    }

    return true;
  }

  bool allSymbolsIsDigits(String str) {
    for (int i = 0; i < str.length; i++) {
      if (!isDigit(str, i)) {
        return false;
      }
    }

    return true;
  }

  bool isDigit(String s, int idx) =>
      "0".compareTo(s[idx]) <= 0 && "9".compareTo(s[idx]) >= 0;
}
