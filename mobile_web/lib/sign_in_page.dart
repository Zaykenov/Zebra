import 'package:email_validator/email_validator.dart';
import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/registration_page.dart';
import 'package:mobile_web/sign_in_validation_page.dart';
import 'package:platform_device_id/platform_device_id.dart';

class SignInPage extends StatefulWidget {
  const SignInPage({Key? key}) : super(key: key);

  @override
  State<SignInPage> createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> {
  final _form = GlobalKey<FormState>();

  final TextEditingController _userEmailController = TextEditingController();

  bool showSentCodeField = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: SingleChildScrollView(
      padding: const EdgeInsets.symmetric(vertical: 100),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          const Image(image: AssetImage('assets/logo.png')),
          const Padding(
            padding: EdgeInsets.all(20.0),
            child: Column(
              children: [
                Align(
                    alignment: Alignment.center,
                    child: Text(
                      "Авторизация",
                      style:
                          TextStyle(fontWeight: FontWeight.bold, fontSize: 20),
                    ))
              ],
            ),
          ),
          Form(
            key: _form,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 50),
              child: Column(
                children: [
                  TextFormField(
                    autofocus: true,
                    controller: _userEmailController,
                    validator: (value) => EmailValidator.validate(value ?? "")
                        ? null
                        : "Введите правильный email",
                    keyboardType: TextInputType.emailAddress,
                    decoration: const InputDecoration(labelText: 'Email'),
                  ),
                ],
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
                String? deviceId = await PlatformDeviceId.getDeviceId;

                SmartDialog.showLoading(msg: "Отправляем вам код :)");

                var status = await ApiCaller().trySignIn(
                    TrySignInModel(_userEmailController.text, deviceId!));
                await SmartDialog.dismiss();

                if (status == TrySignInResStatus.UserNotExists) {
                  SmartDialog.showToast(
                      "Вы еще не зарегистрированы. Попробуйте пройти регистрацию");
                  await Navigator.pushAndRemoveUntil(
                      _form.currentContext!,
                      MaterialPageRoute(
                          builder: (context) => const RegistrationPage()),
                      (route) => false);
                } else {
                  if (status == TrySignInResStatus.AlreadySentToEmail) {
                    SmartDialog.showToast(
                        "Код валидации уже был отправлен ранее");
                  }
                  await Navigator.push(
                      _form.currentContext!,
                      MaterialPageRoute(
                          builder: (context) =>
                              SignInValidationPage(_userEmailController.text)));
                }
              }
            },
            child: const Text(
              "Отправить код",
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
          )
        ],
      ),
    ));
  }
}
