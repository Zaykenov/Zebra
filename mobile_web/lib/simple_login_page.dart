import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:mobile_web/qr_page.dart';
import 'package:form_validator/form_validator.dart';
import 'package:platform_device_id/platform_device_id.dart';

class SimpleLoginPage extends StatefulWidget {
  SimpleLoginPage({Key? key}) : super(key: key);
  final TextEditingController _userNameController = TextEditingController();

  @override
  State<SimpleLoginPage> createState() => _SimpleLoginPageState();
}

class _SimpleLoginPageState extends State<SimpleLoginPage> {
  final _form = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
      child: SingleChildScrollView(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Image(image: AssetImage('assets/logo.png')),
            const Padding(
              padding: EdgeInsets.all(20.0),
              child: Column(
                children: [
                  // Align(
                  //     alignment: Alignment.centerLeft,
                  //     child: Text(
                  //       "Анонимная регистрация",
                  //       style: TextStyle(fontWeight: FontWeight.bold),
                  //     )),
                  Align(
                      alignment: Alignment.centerLeft,
                      child:
                          Text("Без регистрации процент скидки составит 5%")),
                ],
              ),
            ),
            Form(
              key: _form,
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: TextFormField(
                  controller: widget._userNameController,
                  validator:
                      ValidationBuilder().minLength(3).maxLength(50).build(),
                  decoration: const InputDecoration(
                      labelText: 'Пожалуйста представьтесь'),
                ),
              ),
            ),
            const SizedBox(height: 30),
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                foregroundColor: Colors.white,
                backgroundColor: const Color(0xff3EB2B2),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(25),
                ),
                elevation: 15.0,
                padding:
                    const EdgeInsets.symmetric(horizontal: 100, vertical: 20),
              ),
              onPressed: () async {
                if (_form.currentState != null &&
                    _form.currentState!.validate()) {
                  SmartDialog.showLoading(msg: "Осталось совсем чуть-чуть ☕");
                  var userName = widget._userNameController.text;
                  String? deviceId = await PlatformDeviceId.getDeviceId;

                  var createIncognitoUserRes = await ApiCaller()
                      .createIncognitoUser(
                          CreateIncognitoUserReq(userName, deviceId!));
                  var userEmail = "";

                  await SessionHelper().userSignIn(
                      userEmail, userName, createIncognitoUserRes.clientId!);
                  await SmartDialog.dismiss();
                  if (mounted) {
                    await Navigator.pushAndRemoveUntil(
                        context,
                        MaterialPageRoute(builder: (context) => const QrPage()),
                        (route) => false);
                  }
                }
              },
              child: const Text(
                "Продолжить",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            )
          ],
        ),
      ),
    ));
  }
}
