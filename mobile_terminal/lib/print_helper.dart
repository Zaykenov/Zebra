import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';
import 'dart:ui' as ui;
import 'package:flutter_esc_pos_utils/flutter_esc_pos_utils.dart';

import 'package:flutter/material.dart';
import 'package:flutter_pos_printer_platform/flutter_pos_printer_platform.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:intl/intl.dart';
import 'package:mobile_terminal/check_models/check_model.dart';

import 'package:image/image.dart';
import 'package:image/image.dart' as img;
import 'package:barcode/barcode.dart' as barcode;

class PrintHelper {
  Future<void> printWiFiReceipt(WiFiPrinter printer, CheckModel model) async {
    var bytes = await getCheckContentToPrint(model);

    var tcpPrinterInput = TcpPrinterInput(
        ipAddress: printer.ipAddress!,
        port: int.parse(printer.port!),
        timeout: Duration(seconds: int.parse(printer.timeoutInSec!)));

    await PrinterManager.instance.tcpPrinterConnector.connect(tcpPrinterInput);

    var print2x = bytes + bytes;
    await PrinterManager.instance.tcpPrinterConnector.send(print2x);
  }

  Future<void> printUsbReceipt(UsbPrinter printer, CheckModel model,
      {int count = 1}) async {
    var bytes = await getCheckContentToPrint(model);

    await PrinterManager.instance.connect(
        type: PrinterType.usb,
        model: UsbPrinterInput(
            name: printer.deviceName,
            productId: printer.productId,
            vendorId: printer.vendorId));

    while (count > 0) {
      await PrinterManager.instance.send(type: PrinterType.usb, bytes: bytes);
      count--;
    }
  }

  Future<List<int>> getCheckContentToPrint(CheckModel model) async {
    List<int> bytes = [];
    final profile = await CapabilityProfile.load();
    final generator =
        Generator(PaperSize.mm80, profile, codec: const Utf8Codec());

    bytes += generator.image(await getStringAsImage("ZEBRA CRM",
        alig: TextAlign.center, fWeight: ui.FontWeight.bold, fontSize: 28));
    bytes += generator.feed(1);
    bytes += generator.image(await getStringAsImage("Чек №  ${model.id}"));
    bytes += generator.image(await getStringAsImage(
        "Напечатан   ${DateFormat('dd.MM.yyyy HH:mm').format(DateTime.now())}"));

    if (model.comment != null && model.comment!.isNotEmpty) {
      bytes += generator
          .image(await getStringAsImage("Комментарий:   ${model.comment}"));
    }

    bytes += generator.hr();

    List<CheckItem> itemsInCheckList = [];
    for (var x in model.techCartCheck!) {
      var modificators = x.getModificators();
      var dopItems = modificators.map((m) => "${m.name}").toList();
      itemsInCheckList
          .add(CheckItem(x.name!, dopItems, x.comments, x.price!, x.quantity!));
    }
    for (var x in model.tovarCheck!) {
      itemsInCheckList.add(
          CheckItem(x.name!, <String>[], x.comments, x.price!, x.quantity!));
    }

    for (var i = 0; i < itemsInCheckList.length; i++) {
      var x = itemsInCheckList[i];
      bytes += await getCheckItemAsImage(i + 1, x, generator);
    }

    bytes += generator.hr();

    if ((model.discountPercent ?? 0) > 0) {
      bytes += generator.image(await getStringAsImage(
          "Скидка (${model.discountPercent}%): ${NumberFormat.simpleCurrency(name: "KZT").format(model.discount)}"));
    }
    bytes += generator.image(await getStringAsImage(
        "ИТОГО  ${NumberFormat.simpleCurrency(name: "KZT").format(model.sum)}",
        alig: ui.TextAlign.left,
        fWeight: ui.FontWeight.bold,
        fontSize: 25));

    bytes += generator.feed(2);
    final dm = barcode.Barcode.qrCode();
    final svg = dm.toSvg(model.tisCheckUrl!, width: 10, height: 10);
    bytes += generator.image(decodeImage(await svgToPng(svg))!);
    bytes += generator.image(await getStringAsImage(
        "Для просмотра фискального чека отсканируйте QR",
        alig: ui.TextAlign.center,
        fontSize: 18));

    bytes += generator.cut();
    return bytes;
  }

  Future<Uint8List> svgToPng(String svgString) async {
    DrawableRoot svgDrawableRoot = await svg.fromSvgString(svgString, "");

    // to have a nice rendering it is important to have the exact original height and width,
    // the easier way to retrieve it is directly from the svg string
    // but be careful, this is an ugly fix for a flutter_svg problem that works
    // with my images
    int originalHeight = 10;
    int originalWidth = 10;

    // toPicture() and toImage() don't seem to be pixel ratio aware, so we calculate the actual sizes here
    double devicePixelRatio = 15;

    double width = originalHeight *
        devicePixelRatio; // where 32 is your SVG's original width
    double height = originalWidth * devicePixelRatio; // same thing

    // Convert to ui.Picture
    final picture = svgDrawableRoot.toPicture(size: Size(width, height));

    // Convert to ui.Image. toImage() takes width and height as parameters
    // you need to find the best size to suit your needs and take into account the screen DPI
    final image = await picture.toImage(width.toInt(), height.toInt());
    var bytes = await image.toByteData(format: ui.ImageByteFormat.png);

    return bytes!.buffer.asUint8List();
  }
  // String prepare(String input) {
  //   if (input.length < 18) {
  //     return input;
  //   }

  // }

  Future<img.Image> getStringAsImage(String text,
      {TextAlign alig = TextAlign.left,
      ui.FontWeight fWeight = ui.FontWeight.normal,
      double fontSize = 25}) async {
    return decodeImage(await _generateImageFromString(text, alig,
        fWeight: fWeight, fontSize: fontSize))!;
  }

  Future<List<int>> getCheckItemAsImage(
      int indexInCheck, CheckItem item, Generator generator) async {
    List<int> bytes = [];

    final formatCurrency = NumberFormat.simpleCurrency(name: "KZT");
    bytes += generator.image(await getStringAsImage("◆ ${item.name}",
        fWeight: ui.FontWeight.bold, fontSize: 28));

    for (var dopItem in item.dopItems) {
      bytes += generator
          .image(await getStringAsImage("\t ◆ $dopItem", fontSize: 22));
    }
    if (item.comment != null && item.comment!.isNotEmpty) {
      bytes += generator
          .image(await getStringAsImage("\t ◇ ${item.comment}", fontSize: 22));
    }

    bytes += generator
        .image(await getStringAsImage("кол-во: ${item.count}", fontSize: 28));
    bytes += generator.image(await getStringAsImage(
        "x ${formatCurrency.format(item.price)}   =  ${formatCurrency.format(item.price * item.count)}",
        fontSize: 20));
    return bytes;
  }

  Future<Uint8List> _generateImageFromString(String text, ui.TextAlign align,
      {ui.FontWeight fWeight = ui.FontWeight.normal,
      double fontSize = 25}) async {
    ui.PictureRecorder recorder = ui.PictureRecorder();
    Canvas canvas = Canvas(
        recorder,
        Rect.fromCenter(
          center: const Offset(0, 0),
          width: 600,
          height: 400, // cheated value, will will clip it later...
        ));
    TextSpan span = TextSpan(
      style: TextStyle(
          color: Colors.black,
          fontSize: fontSize,
          leadingDistribution: ui.TextLeadingDistribution.proportional,
          fontWeight: fWeight,
          wordSpacing: 0),
      text: text,
    );
    TextPainter tp = TextPainter(
        text: span,
        maxLines: 3,
        textAlign: align,
        textDirection: ui.TextDirection.ltr);
    tp.layout(minWidth: 600, maxWidth: 600);
    tp.paint(canvas, const Offset(0.0, 0.0));
    var picture = recorder.endRecording();
    final pngBytes = await picture.toImage(
      tp.size.width.toInt(),
      tp.size.height.toInt() - 3, // decrease padding
    );
    final byteData = await pngBytes.toByteData(format: ui.ImageByteFormat.png);

    //var bse64 = base64Encode(byteData!.buffer.asUint8List());
    return byteData!.buffer.asUint8List();
  }

  Future<void> printReceipt(String url) async {
    // await Printing.layoutPdf(
    //     onLayout: (PdfPageFormat format) async => await _getPdfFromUrl(url),
    //     format: PdfPageFormat.roll80);
  }
}

class DefaultPrinter {
  WiFiPrinter? wiFiPrinter;
  UsbPrinter? usbPrinter;

  DefaultPrinter({
    this.wiFiPrinter,
    this.usbPrinter,
  });

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['wiFiPrinter'] = wiFiPrinter == null ? null : wiFiPrinter!.toJson();
    data['usbPrinter'] = usbPrinter == null ? null : usbPrinter!.toJson();
    return data;
  }

  DefaultPrinter.fromJson(Map<String, dynamic> json) {
    wiFiPrinter = json['wiFiPrinter'] == null
        ? null
        : WiFiPrinter.fromJson(json['wiFiPrinter']);
    usbPrinter = json['usbPrinter'] == null
        ? null
        : UsbPrinter.fromJson(json['usbPrinter']);
  }
}

class WiFiPrinter {
  String? ipAddress;
  String? port;
  String? timeoutInSec;

  WiFiPrinter({
    this.ipAddress,
    this.port,
    this.timeoutInSec,
  });

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['ipAddress'] = ipAddress;
    data['port'] = port;
    data['timeoutInSec'] = timeoutInSec;
    return data;
  }

  WiFiPrinter.fromJson(Map<String, dynamic> json) {
    ipAddress = json['ipAddress'];
    port = json['port'];
    timeoutInSec = json['timeoutInSec'];
  }
}

class UsbPrinter {
  int? id;
  String? deviceName;
  String? address;
  String? port;
  String? vendorId;
  String? productId;

  UsbPrinter({
    this.deviceName,
    this.address,
    this.port,
    this.vendorId,
    this.productId,
  });

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['id'] = id;
    data['deviceName'] = deviceName;
    data['address'] = address;
    data['port'] = port;
    data['vendorId'] = vendorId;
    data['productId'] = productId;
    return data;
  }

  UsbPrinter.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    deviceName = json['deviceName'];
    address = json['address'];
    port = json['port'];
    vendorId = json['vendorId'];
    productId = json['productId'];
  }
}
