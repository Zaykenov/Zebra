import 'dart:convert';

import 'package:mobile_terminal/print_helper.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPrefHelper {
  Future<DefaultPrinter?> getSavedPrinterOrNull() async {
    final prefs = await SharedPreferences.getInstance();
    var json = prefs.getString('printer_config1');

    if (json == null || json.isEmpty) {
      return null;
    }

    return DefaultPrinter.fromJson(jsonDecode(json));
  }

  Future<void> savePrinter(DefaultPrinter p) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('printer_config1', jsonEncode(p.toJson()));
  }
}
