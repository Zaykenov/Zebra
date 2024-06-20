import 'package:email_validator/email_validator.dart';
import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:form_validator/form_validator.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/registration_validation_page.dart';
import 'package:platform_device_id/platform_device_id.dart';

import 'sign_in_page.dart';

class RegistrationPage extends StatefulWidget {
  const RegistrationPage({Key? key}) : super(key: key);

  @override
  State<RegistrationPage> createState() => _RegistrationPageState();
}

class _RegistrationPageState extends State<RegistrationPage> {
  final _form = GlobalKey<FormState>();

  final TextEditingController _userEmailController = TextEditingController();
  final TextEditingController _userNameController = TextEditingController();
  DateTime? birthDate;

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
                    alignment: Alignment.centerLeft,
                    child: Text(
                      "Регистрация по email",
                      style: TextStyle(fontWeight: FontWeight.bold),
                    )),
                Align(
                    alignment: Alignment.centerLeft,
                    child: Text(
                        "Она позволяет вам получить 10% при любом вашем заказе")),
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
                    controller: _userEmailController,
                    validator: (value) => EmailValidator.validate(value ?? "")
                        ? null
                        : "Введите правильный email",
                    keyboardType: TextInputType.emailAddress,
                    decoration: const InputDecoration(labelText: 'Email'),
                  ),
                  TextFormField(
                    controller: _userNameController,
                    validator:
                        ValidationBuilder().minLength(3).maxLength(50).build(),
                    decoration: const InputDecoration(labelText: 'Ваше имя'),
                  ),
                  // DateTimeFormField(
                  //   decoration: const InputDecoration(
                  //     hintStyle: TextStyle(color: Colors.black45),
                  //     errorStyle:
                  //         TextStyle(color: Color.fromARGB(255, 179, 38, 38)),
                  //     //border: OutlineInputBorder(),
                  //     suffixIcon: Icon(Icons.event_note),
                  //     labelText: 'Дата рождения*',
                  //   ),

                  //   lastDate: DateTime.now(),
                  //   dateFormat: DateFormat("dd.MM.yyyy"),
                  //   firstDate: DateTime(1900, 01, 01),
                  //   initialDatePickerMode: DatePickerMode.year,
                  //   mode: DateTimeFieldPickerMode.date,
                  //   //autovalidateMode: AutovalidateMode.always,
                  //   validator: (e) {
                  //     if (e == null) {
                  //       return "Пустая дата рождения";
                  //     }
                  //     return null;
                  //   },
                  //   onDateSelected: (DateTime value) {
                  //     setState(() {
                  //       birthDate = value;
                  //     });
                  //   },
                  // ),
                  // const Text(
                  //     "* дата рождения нужна нам что бы успеть подготовить вам подарок")
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

                SmartDialog.showLoading(msg: "Регистрируем вас)");

                var email = _userEmailController.text.trim();
                var userName = _userNameController.text.trim();
                var status = await ApiCaller().tryCreateUser(
                    TryCreateUserReq(email, userName, birthDate, deviceId!));
                await SmartDialog.dismiss();

                if (status == TryCreateUserResStatus.AlreadyRegistered) {
                  SmartDialog.showToast(
                      "Вы уже зарегистрированы. Попробуйте войти");
                  await Navigator.pushAndRemoveUntil(
                      _form.currentContext!,
                      MaterialPageRoute(
                          builder: (context) => const SignInPage()),
                      (route) => false);
                } else {
                  if (status == TryCreateUserResStatus.AlreadyRegistered) {
                    SmartDialog.showToast("Код валидации уже был отправлен");
                  }
                  await Navigator.push(
                      _form.currentContext!,
                      MaterialPageRoute(
                          builder: (context) =>
                              RegistrationValidationPage(email, userName)));
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
    ));
  }
}
