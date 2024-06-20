import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:liquid_pull_to_refresh/liquid_pull_to_refresh.dart';
import 'package:mobile_web/about_page.dart';
import 'package:mobile_web/feedback_page.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/render_helper.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:mobile_web/main.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:shimmer/shimmer.dart';

import 'map/live_location.dart';

class QrPage extends StatefulWidget {
  const QrPage({Key? key}) : super(key: key);

  @override
  State<QrPage> createState() => _QrPageState();
}

class _QrPageState extends State<QrPage> {
  final GlobalKey<LiquidPullToRefreshState> _refreshIndicatorKey =
      GlobalKey<LiquidPullToRefreshState>();
  String? userName;
  String? userAccessCode;
  double? zebraCoinBalance;
  double? discountInPercent;
  List<OrderAndFeedback>? currentOrders;

  StreamSubscription? _updateCurrentOrdersSubs;
  StreamSubscription? _updateQrSubs;
  StreamSubscription? _updateCoinsSubs;

  bool? needToRenderDemoBunner;

  void startCurrentOrdersTimer() async {
    await cancelCurrentOrdersTimer();
    _updateCurrentOrdersSubs =
        Stream.periodic(const Duration(seconds: 5)).listen((event) {
      if (mounted) {
        loadCurrentOrders();
      }
    });
  }

  Future<void> cancelCurrentOrdersTimer() async {
    await _updateCurrentOrdersSubs?.cancel();
  }

  void startCoinsTimer() async {
    await cancelCoinsTimer();
    _updateQrSubs =
        Stream.periodic(const Duration(seconds: 30)).listen((event) {
      if (mounted) {
        loadZebraCoinsAndDiscount();
      }
    });
  }

  Future<void> cancelCoinsTimer() async {
    await _updateCoinsSubs?.cancel();
  }

  void startQrTimer() async {
    await cancelQrTimer();
    _updateQrSubs =
        Stream.periodic(const Duration(seconds: 30)).listen((event) {
      if (mounted) {
        loadQr();
      }
    });
  }

  Future<void> cancelQrTimer() async {
    await _updateQrSubs?.cancel();
  }

  @override
  void initState() {
    super.initState();
    startLoadUserData();
    startCurrentOrdersTimer();
    startQrTimer();
    startCoinsTimer();
  }

  @override
  void dispose() {
    cancelCurrentOrdersTimer();
    cancelQrTimer();
    cancelCoinsTimer();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    initializeDateFormatting('ru', null);
    return Scaffold(
      backgroundColor: const Color(0xFFF8F7F7),
      body: LiquidPullToRefresh(
        key: _refreshIndicatorKey,
        onRefresh: () async {
          await startLoadUserData();
        },
        animSpeedFactor: 1,
        springAnimationDurationInMilliseconds: 200,
        color: const Color(0xFF3EB2B2),
        child: ListView(children: [
          Padding(
            padding: const EdgeInsets.all(10.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: <Widget>[
                getHeaderWidget(userName),
                const SizedBox(height: 10),
                Row(
                  children: [
                    // Expanded(
                    //   child:
                    //       getZebraCoinBalanceOrShimmerWidget(zebraCoinBalance),
                    // ),
                    // const SizedBox(width: 10),
                    Expanded(
                      child: getDiscountInPercentOrShimmerWidget(
                          discountInPercent),
                    ),
                  ],
                ),
                const SizedBox(
                  height: 20,
                ),
                SizedBox(
                  width: double.infinity,
                  height: 110,
                  child: Material(
                    elevation:
                        4, // Adjust the elevation value as needed for the desired shadow intensity
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(30),
                    ),
                    color: Colors.white,
                    child: getQrOrShimmerWidget(userAccessCode),
                  ),
                ),
                const SizedBox(height: 10),
                const Padding(
                  padding: EdgeInsets.all(8.0),
                  child: Text(
                    "Заказы:",
                    style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 20,
                        color: Color(0xff717171)),
                  ),
                ),
                getCurrentOrdersOrShimmerWidget(currentOrders),
              ],
            ),
          ),
        ]),
      ),
    );
  }

  Widget getHeaderWidget(String? userName) {
    return Row(
      children: [
        IconButton(
            onPressed: () async {
              SmartDialog.show(builder: (context) {
                return Container(
                  height: 200,
                  width: 300,
                  decoration: BoxDecoration(
                    color: const Color(0xfffafafa),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  alignment: Alignment.center,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(top: 10),
                        child: Text('Пользователь: ${userName ?? "-"}',
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                                fontWeight: FontWeight.bold, fontSize: 18)),
                      ),
                      const SizedBox(
                        height: 10,
                      ),
                      TextButton(
                        style: ButtonStyle(
                          foregroundColor:
                              MaterialStateProperty.all<Color>(Colors.blue),
                        ),
                        onPressed: () async {
                          SmartDialog.dismiss();
                          await Navigator.push(
                              _refreshIndicatorKey.currentContext!,
                              MaterialPageRoute(
                                  builder: (context) => const AboutPage()));
                        },
                        child: const Text('О программе',
                            style: TextStyle(color: Colors.black)),
                      ),
                      TextButton(
                        style: ButtonStyle(
                          foregroundColor:
                              MaterialStateProperty.all<Color>(Colors.blue),
                        ),
                        onPressed: () async {
                          await SessionHelper().userSignOut().then(
                              (value) => RestartWidget.restartApp(context));
                        },
                        child: const Text('Выйти',
                            style: TextStyle(color: Colors.red)),
                      )
                    ],
                  ),
                );
              });
            },
            icon: const Icon(Icons.account_circle_outlined)),
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Text(userName ?? "-", style: const TextStyle(fontSize: 16)),
        ),
        const Spacer(),
        (needToRenderDemoBunner ?? false)
            ? const Text("_______",
                style: TextStyle(
                    color: Colors.white,
                    backgroundColor: Color.fromARGB(255, 52, 211, 153)))
            : Container(),
        const Spacer(),
        // IconButton(
        //   onPressed: () async {
        //     await Navigator.push(
        //         _refreshIndicatorKey.currentContext!,
        //         MaterialPageRoute(
        //             builder: (context) => const ChooseLocation()));
        //   },
        //   icon: const Icon(Icons.shopping_bag_outlined),
        // ),
        IconButton(
          onPressed: () async {
            if (await Permission.locationWhenInUse.request().isGranted) {
              await Navigator.push(_refreshIndicatorKey.currentContext!,
                  MaterialPageRoute(builder: (context) => MapScreen()));
            }
          },
          icon: const Icon(Icons.map_outlined),
        )
      ],
    );
  }

  Widget getQrOrShimmerWidget(String? userAccessCode) {
    if (userAccessCode == null) {
      return Shimmer.fromColors(
        baseColor: const Color.fromARGB(255, 252, 250, 250),
        highlightColor: const Color.fromARGB(255, 230, 230, 230),
        child: Container(
          width: 50,
          height: 50,
          color: Colors.white,
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.only(left: 20, top: 10),
          child: Text(
            'Продиктуйте код кассиру',
            style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: Color(0xFFAEAEAE)),
          ),
        ),
        Center(
          child: Text(
            userAccessCode,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 48),
          ),
        )
      ],
    );
  }

  Widget getZebraCoinBalanceOrShimmerWidget(double? zebraCoinBalance) {
    if (userAccessCode == null) {
      return Shimmer.fromColors(
        baseColor: const Color.fromARGB(255, 252, 250, 250),
        highlightColor: const Color.fromARGB(255, 230, 230, 230),
        child: Container(
          width: 50,
          height: 50,
          color: Colors.white,
        ),
      );
    }

    void showZebraCoinDialog() {
      showDialog(
        context: context,
        builder: (BuildContext context) {
          return const SimpleDialog(
            children: [
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 20, vertical: 20),
                child: Text(
                  'Добро пожаловать в мир Zebra Coffee и систему Zebra Coin! Получайте монеты за регистрацию, покупки и комментарии🦓. Обменяйте их на журналы, термокружки, мерч, камри 80, квартиру на левом берегу Жезказгана... Кликните на иконку Zebra Coin и окунитесь в увлекательный мир вознаграждений! ☕️ Данный функционал в настоящее время находится в разработке и скоро будет доступен для вас🫶',
                  style: TextStyle(fontSize: 16),
                ),
              ),
              // SimpleDialogOption(
              //   onPressed: () {
              //     Navigator.pop(context);
              //   },
              //   child: const Text('OK'),
              // ),
            ],
          );
        },
      );
    }

    return GestureDetector(
      onTap: showZebraCoinDialog,
      child: SizedBox(
        width: double.infinity,
        height: 100,
        child: Material(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          color: Colors.white,
          elevation: 2,
          child: InkWell(
            borderRadius: BorderRadius.circular(20),
            onTap: showZebraCoinDialog,
            child: Padding(
              padding: const EdgeInsets.only(left: 15, top: 10),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Ваши зебракоины',
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.bold,
                        color: Color(0xFFAEAEAE)),
                  ),
                  Row(
                    children: [
                      Text(
                        "${zebraCoinBalance?.round() ?? "-"}",
                        style: const TextStyle(
                            color: Color(0xFF3EB2B2),
                            fontWeight: FontWeight.bold,
                            fontSize: 48),
                      ),
                      const SizedBox(
                        width: 5,
                      ),
                      Container(
                        width: 40, // Adjust the width as needed
                        height: 40, // Adjust the height as needed
                        decoration: const BoxDecoration(
                          image: DecorationImage(
                            image: AssetImage(
                                'assets/zebraCoinV2.png'), // Replace with your image asset path
                            fit: BoxFit.cover,
                          ),
                          shape: BoxShape.circle,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget getDiscountInPercentOrShimmerWidget(double? discountInPercent) {
    if (userAccessCode == null) {
      return Shimmer.fromColors(
        baseColor: const Color.fromARGB(255, 252, 250, 250),
        highlightColor: const Color.fromARGB(255, 230, 230, 230),
        child: Container(
          width: double.infinity,
          height: 100,
          color: Colors.white,
        ),
      );
    }

    return SizedBox(
      width: double.infinity,
      height: 100,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(20),
          color: const Color.fromARGB(180, 62, 178, 178),
        ),
        child: Padding(
          padding: const EdgeInsets.only(left: 15, top: 10),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Процент скидки',
                style: TextStyle(
                    fontSize: 15,
                    color: Colors.white,
                    fontWeight: FontWeight.bold),
              ),
              Text(
                "${discountInPercent?.round() ?? "-"}%",
                style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 48),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget getCurrentOrdersOrShimmerWidget(
      List<OrderAndFeedback>? currentOrders) {
    if (currentOrders == null) {
      return Shimmer.fromColors(
        baseColor: const Color.fromARGB(255, 252, 250, 250),
        highlightColor: const Color.fromARGB(255, 230, 230, 230),
        child: Container(
          width: double.infinity,
          height: 100,
          color: Colors.white,
        ),
      );
    }

    if (currentOrders.isEmpty) {
      return const Padding(
        padding: EdgeInsets.all(8.0),
        child: Text("Заказов нет"),
      );
    }

    List<Widget> orderWidgetList = <Widget>[];
    currentOrders.sort((a, b) {
      return -a.order.openedat.compareTo(b.order.openedat);
    });
    for (var orderAndFeedback in currentOrders) {
      orderWidgetList.add(InkWell(
        onTap: () async {
          if (mounted) {
            await Navigator.push(
                context,
                MaterialPageRoute(
                    builder: (context) => FeedbackPage(orderAndFeedback)));
          }
        },
        child: Column(
          children: [
            Stack(
              children: [
                Card(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(30),
                  ),
                  color: const Color.fromARGB(180, 62, 178, 178),
                  child: Column(
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(top: 10, left: 20),
                        child: RenderHelper(
                          check: orderAndFeedback.order,
                          needToRenderStar: true,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(
              height: 10,
            ),
          ],
        ),
      ));
    }

    return Column(
      children: orderWidgetList,
    );
  }

  Future loadAppEnv() async {
    var conf = await SessionHelper().getAppConfig();
    setState(() {
      needToRenderDemoBunner = conf.envName == EnvName.staging;
    });
  }

  Future startLoadUserData() async {
    await loadUserName();
    await loadQr();
    await loadZebraCoinsAndDiscount();
    await loadCurrentOrders();
    await loadAppEnv();
  }

  Future loadUserName() async {
    var userNameInStorage = await SessionHelper().getClientUserName();
    setState(() {
      userName = userNameInStorage;
    });
  }

  Future loadQr() async {
    var res = await ApiCaller().genQrCode(await SessionHelper().getClientId());
    if (res.status == TryGenerateQrResStatus.Ok) {
      setState(() {
        userAccessCode = res.qrContent;
      });
    } else {
      SmartDialog.showToast("Неудалось сгенирировать qr");
    }
  }

  Future loadZebraCoinsAndDiscount() async {
    var res =
        await ApiCaller().getClientInfo(await SessionHelper().getClientId());

    if (res.status == TryGetUserInfoResStatus.Ok) {
      setState(() {
        discountInPercent = res.userInfo!.discount;
        zebraCoinBalance = res.userInfo!.zebraCoinBalance;
      });
    } else {
      SmartDialog.showToast("Неудалось подгрузить скидку");
    }
  }

  Future loadCurrentOrders() async {
    var res = await ApiCaller()
        .getCurrentOrdersRes(await SessionHelper().getClientId());

    if (res.status == TryGetLastOrdersResStatus.Ok) {
      var ordersAndFeedbacks = <OrderAndFeedback>[];
      for (var order in res.orders!) {
        var feedback = await SessionHelper().getFeedbackForOrder(order.id);
        ordersAndFeedbacks.add(OrderAndFeedback(order, feedback));
      }

      setState(() {
        currentOrders = ordersAndFeedbacks;
      });
    } else {
      SmartDialog.showToast("Неудалось подгрузить заказы");
    }
  }
}
