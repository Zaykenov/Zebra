import 'package:flutter/material.dart';
import 'package:mobile_web/registration_page.dart';
import 'package:mobile_web/simple_login_page.dart';

import 'helpers/session_helper.dart';
import 'sign_in_page.dart';

class ChooseSignInTypePage extends StatefulWidget {
  const ChooseSignInTypePage({Key? key}) : super(key: key);

  @override
  State<ChooseSignInTypePage> createState() => _ChooseSignInTypePageState();
}

class _ChooseSignInTypePageState extends State<ChooseSignInTypePage> {
  EnvName? envName;

  Future loadAppEnv() async {
    var conf = await SessionHelper().getAppConfig();
    setState(() {
      envName = conf.envName;
    });
  }

  @override
  void initState() {
    super.initState();
    loadAppEnv();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          const Text("стань самым любимым\nпосетителем",
              textAlign: TextAlign.center, style: TextStyle(fontSize: 22)),
          const Image(image: AssetImage('assets/logo.png')),
          const SizedBox(height: 30),
          SizedBox(
            width: 280,
            height: 50,
            child: ElevatedButton(
              style: ElevatedButton.styleFrom(
                foregroundColor: Colors.white,
                backgroundColor: const Color.fromARGB(255, 62, 178, 178),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(25),
                ),
                elevation: 0,
              ),
              onPressed: () async {
                await Navigator.push(
                    context,
                    MaterialPageRoute(
                        builder: (context) => const SignInPage()));
              },
              child: const Text(
                "Войти",
                style: TextStyle(fontWeight: FontWeight.bold, fontSize: 18),
              ),
            ),
          ),
          const SizedBox(
            height: 10,
          ),
          SizedBox(
            width: 280,
            height: 50,
            child: ElevatedButton(
              style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(25),
                      side: const BorderSide(color: Colors.black26)),
                  elevation: 0),
              onPressed: () async {
                await Navigator.push(
                    context,
                    MaterialPageRoute(
                        builder: (context) => const RegistrationPage()));
              },
              child: const Text(
                "Регистрация",
                style: TextStyle(color: Colors.black45, fontSize: 18),
              ),
            ),
          ),
          TextButton(
              style: ButtonStyle(
                overlayColor: MaterialStateProperty.all<Color>(
                  const Color.fromARGB(38, 62, 178, 178),
                ),
              ),
              onPressed: () async {
                await Navigator.push(context,
                    MaterialPageRoute(builder: (context) => SimpleLoginPage()));
              },
              child: const Text(
                "продолжить без регистрации",
                style: TextStyle(color: Colors.black38),
              ))
        ],
      ),
    ));
  }
}
