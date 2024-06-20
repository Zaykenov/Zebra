import 'dart:async';
import 'dart:convert';

import 'package:buttons_tabbar/buttons_tabbar.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_pos_printer_platform/flutter_pos_printer_platform.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_terminal/check_models/check_model.dart';
import 'package:mobile_terminal/print_helper.dart' as print_helper;

import 'cache/shared_pref_helper.dart';
import 'input_formatters/ip_input_formatter.dart';

class PrinterSettings extends StatefulWidget {
  const PrinterSettings({Key? key}) : super(key: key);

  @override
  PrinterSettingsState createState() => PrinterSettingsState();
}

class PrinterSettingsState extends State<PrinterSettings> {
  var usbDevices = <print_helper.UsbPrinter>[];

  var wifiPrinterIpTextController = TextEditingController();
  var wifiPrinterPortTextController = TextEditingController();
  var wifiPrinterTimeoutInSecTextController = TextEditingController();

  print_helper.DefaultPrinter? _savedPrinter;

  var printerManager = PrinterManager.instance;

  final String testReceipJson = """
{
    "tisCheckUrl": "https://app.kassa.wipon.kz/links/check/eace60c8-4e1b-4708-9c87-a161e7db3cdd",
    "closed_at": "2022-12-30T19:21:43.335084618+06:00",
    "comment": "Побыстрее",
    "cost": 840.94104,
    "discount": 158,
    "discount_percent": 10,
    "feedback": null,
    "id": 2360,
    "opened_at": "2022-12-30T19:20:56.430966+06:00",
    "payment": "картой",
    "status": "closed",
    "sum": 1580,
    "techCartCheck": [
        {
           "id": 39666,
            "check_id": 10645,
            "tech_cart_id": 257,
            "name": "Голубая латте матча",
            "quantity": 1,
            "cost": 639.9114,
            "price": 1150,
            "discount": 0,
            "modificators": "[{\\"name\\":\\"Овсяное молоко\\",\\"price\\":300}]",
            "comments": "Слаще"
        }
    ],
    "tovarCheck": [
        {
            "check_id": 2360,
            "comments": "Сладкий как Сабина:)",
            "cost": 333.11,
            "discount": 0,
            "id": 8135,
            "modifications": "",
            "name": "Пончики",
            "price": 590,
            "quantity": 1,
            "tovar_id": 24
        }
    ],
    "user_id": 100,
    "worker": 7
}
""";

  @override
  void initState() {
    super.initState();
    _scanUsbPrinters();
    _setDefaultPrinter();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: DefaultTabController(
      length: 2,
      child: Column(
        children: <Widget>[
          const SizedBox(height: 40),
          Column(
            children: [
              Text(
                "Сейчас используется принтер: ${_savedPrinter == null ? "не указан" : _getPrinterStr(_savedPrinter!)}",
                style:
                    const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              )
            ],
          ),
          const SizedBox(height: 10),
          ButtonsTabBar(
            backgroundColor: Colors.blue,
            //contentPadding: EdgeInsets.all(20),
            radius: 1,
            tabs: const <Widget>[
              Tab(
                icon: Icon(Icons.usb),
                text: "USB",
              ),
              Tab(
                icon: Icon(Icons.wifi),
                text: "WI-FI",
              ),
            ],
          ),
          const Divider(color: Colors.black),
          Expanded(
            child: TabBarView(
              children: <Widget>[
                Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: <Widget>[
                      TextButton(
                        child: const Text("Обновить список"),
                        onPressed: () {
                          _scanUsbPrinters();
                        },
                      ),
                      usbDevices.isEmpty
                          ? const Center(child: Text("Принтера не найдены"))
                          : Expanded(
                              child: ListView(
                                  padding: const EdgeInsets.all(30),
                                  children: usbDevices
                                      .map((device) => GestureDetector(
                                            onTap: () {
                                              SmartDialog.show(
                                                  builder: (context) {
                                                return Container(
                                                  height: 200,
                                                  width: 300,
                                                  decoration: BoxDecoration(
                                                    color: Colors.white,
                                                    borderRadius:
                                                        BorderRadius.circular(
                                                            10),
                                                  ),
                                                  alignment: Alignment.center,
                                                  child: Padding(
                                                    padding:
                                                        const EdgeInsets.all(
                                                            8.0),
                                                    child: Column(
                                                      mainAxisAlignment:
                                                          MainAxisAlignment
                                                              .center,
                                                      children: [
                                                        Text(
                                                          'Название принтера: ${device.deviceName ?? "-"}',
                                                        ),
                                                        const SizedBox(
                                                            height: 15),
                                                        ElevatedButton(
                                                          onPressed: () async {
                                                            var newPrinter = print_helper
                                                                .DefaultPrinter(
                                                                    wiFiPrinter:
                                                                        null,
                                                                    usbPrinter:
                                                                        device);
                                                            await SharedPrefHelper()
                                                                .savePrinter(
                                                                    newPrinter);

                                                            setState(() {
                                                              _savedPrinter =
                                                                  newPrinter;
                                                            });

                                                            SmartDialog
                                                                .dismiss();
                                                          },
                                                          style: ElevatedButton.styleFrom(
                                                              elevation: 12.0,
                                                              textStyle:
                                                                  const TextStyle(
                                                                      color: Colors
                                                                          .white)),
                                                          child: const Text(
                                                              'Использовать этот принтер'),
                                                        ),
                                                        const SizedBox(
                                                            height: 10),
                                                        TextButton(
                                                          onPressed: () {
                                                            try {
                                                              SmartDialog
                                                                  .showLoading();
                                                              var model = CheckModel
                                                                  .fromJson(
                                                                      jsonDecode(
                                                                          testReceipJson));

                                                              print_helper
                                                                      .PrintHelper()
                                                                  .printUsbReceipt(
                                                                      device,
                                                                      model,
                                                                      count: 1);
                                                              SmartDialog
                                                                  .dismiss();
                                                            } catch (e) {
                                                              SmartDialog
                                                                  .dismiss();
                                                              SmartDialog.show(
                                                                  builder:
                                                                      (context) {
                                                                return Container(
                                                                  height: 80,
                                                                  width: 180,
                                                                  decoration:
                                                                      BoxDecoration(
                                                                    color: Colors
                                                                        .black,
                                                                    borderRadius:
                                                                        BorderRadius.circular(
                                                                            10),
                                                                  ),
                                                                  alignment:
                                                                      Alignment
                                                                          .center,
                                                                  child: Text(
                                                                      "Ошибка при распечатке: $e",
                                                                      style: const TextStyle(
                                                                          color:
                                                                              Colors.white)),
                                                                );
                                                              });
                                                            }
                                                          },
                                                          style: TextButton.styleFrom(
                                                              textStyle:
                                                                  const TextStyle(
                                                                      color: Colors
                                                                          .white)),
                                                          child: const Text(
                                                              'Тестовая печать'),
                                                        ),
                                                      ],
                                                    ),
                                                  ),
                                                );
                                              });
                                            },
                                            child: Card(
                                              // color: _savedPrinter?.deviceName == device.deviceName
                                              //     ? Colors.amberAccent
                                              //     : Colors.white,

                                              child: Column(
                                                children: [
                                                  Row(
                                                    mainAxisAlignment:
                                                        MainAxisAlignment
                                                            .center,
                                                    children: [
                                                      Column(
                                                        children: [
                                                          Padding(
                                                            padding:
                                                                const EdgeInsets
                                                                    .all(10.0),
                                                            child: Text(device
                                                                    .deviceName ??
                                                                "-"),
                                                          ),
                                                        ],
                                                      ),
                                                    ],
                                                  ),
                                                ],
                                              ),
                                            ),
                                          ))
                                      .toList()),
                            ),
                    ]),
                ListView(
                  //shrinkWrap: true,
                  padding: const EdgeInsets.symmetric(horizontal: 200),
                  children: [
                    SizedBox(
                      width: 255,
                      child: TextFormField(
                        controller: wifiPrinterIpTextController,
                        inputFormatters: [
                          IpAddressInputFilter.intInputFilter(),
                          LengthLimitingTextInputFormatter(15),
                          IpAddressInputFormatter()
                        ],
                        decoration:
                            const InputDecoration(label: Text('Ip адрес')),
                      ),
                    ),
                    SizedBox(
                      width: 255,
                      child: TextFormField(
                        controller: wifiPrinterPortTextController,
                        inputFormatters: [
                          IpAddressInputFilter.intInputFilter(),
                          LengthLimitingTextInputFormatter(6),
                        ],
                        decoration: const InputDecoration(label: Text('Порт')),
                      ),
                    ),
                    SizedBox(
                      width: 255,
                      child: TextFormField(
                        controller: wifiPrinterTimeoutInSecTextController,
                        inputFormatters: [
                          IpAddressInputFilter.intInputFilter(),
                          LengthLimitingTextInputFormatter(10),
                        ],
                        decoration:
                            const InputDecoration(label: Text('Таймаут (сек)')),
                      ),
                    ),
                    TextButton(
                      onPressed: () {
                        try {
                          SmartDialog.showLoading();
                          var model =
                              CheckModel.fromJson(jsonDecode(testReceipJson));

                          var ip = wifiPrinterIpTextController.text;
                          var port = wifiPrinterPortTextController.text;
                          var timeoutInSec =
                              wifiPrinterTimeoutInSecTextController.text;

                          var newPrinter = print_helper.DefaultPrinter(
                              wiFiPrinter: print_helper.WiFiPrinter(
                                  ipAddress: ip,
                                  port: port,
                                  timeoutInSec: timeoutInSec),
                              usbPrinter: null);
                          print_helper.PrintHelper()
                              .printWiFiReceipt(newPrinter.wiFiPrinter!, model);
                          SmartDialog.dismiss();
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
                      },
                      style: TextButton.styleFrom(
                          textStyle: const TextStyle(color: Colors.white)),
                      child: const Text('Тестовая печать'),
                    ),
                    const SizedBox(height: 10),
                    ElevatedButton(
                      style: ButtonStyle(
                          backgroundColor:
                              MaterialStateProperty.all<Color>(Colors.blue)),
                      onPressed: () async {
                        var ip = wifiPrinterIpTextController.text;
                        var port = wifiPrinterPortTextController.text;
                        var timeoutInSec =
                            wifiPrinterTimeoutInSecTextController.text;

                        var newPrinter = print_helper.DefaultPrinter(
                            wiFiPrinter: print_helper.WiFiPrinter(
                                ipAddress: ip,
                                port: port,
                                timeoutInSec: timeoutInSec),
                            usbPrinter: null);
                        await SharedPrefHelper().savePrinter(newPrinter);

                        setState(() {
                          _savedPrinter = newPrinter;
                        });

                        SmartDialog.showToast("Сохранено");
                      },
                      child: const Text("Сохранить"),
                    ),
                  ],
                )
              ],
            ),
          ),
        ],
      ),
    ));
  }

  Future<void> _scanUsbPrinters() async {
    usbDevices.clear();
  }

  Future<void> _setDefaultPrinter() async {
    var defaultPrinter = await SharedPrefHelper().getSavedPrinterOrNull();

    await _fillWifiInputsIfNeed(defaultPrinter);
    setState(() {
      _savedPrinter = defaultPrinter;
    });
  }

  Future<void> _fillWifiInputsIfNeed(
      print_helper.DefaultPrinter? printer) async {
    if (printer?.wiFiPrinter != null) {
      wifiPrinterIpTextController.text = printer!.wiFiPrinter!.ipAddress!;
      wifiPrinterPortTextController.text = printer.wiFiPrinter!.port!;
      wifiPrinterTimeoutInSecTextController.text =
          printer.wiFiPrinter!.timeoutInSec!;
    }
  }

  String _getPrinterStr(print_helper.DefaultPrinter printer) {
    if (printer.usbPrinter != null) {
      return "Тип: USB; Наименование: ${printer.usbPrinter!.deviceName}";
    }

    if (printer.wiFiPrinter != null) {
      return "Тип: WiFi; Ip адрес: ${printer.wiFiPrinter!.ipAddress}; порт: ${printer.wiFiPrinter!.port}";
    }

    return "не распознан";
  }
}
