import 'package:flutter/material.dart';

class RemoveAppPage extends StatefulWidget {
  const RemoveAppPage({Key? key}) : super(key: key);

  @override
  State<RemoveAppPage> createState() => _RemoveAppPageState();
}

class _RemoveAppPageState extends State<RemoveAppPage> {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
            primaryColor: const Color(0xff3EB2B2),
            colorScheme: const ColorScheme.light(primary: Color(0xff3EB2B2)),
            fontFamily: 'Kameron',
            useMaterial3: true),
        home: const Scaffold(
            body: SingleChildScrollView(
          padding: EdgeInsets.symmetric(vertical: 100),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: <Widget>[
              Image(image: AssetImage('assets/logo.png')),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 15, vertical: 30),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      "Мы готовим для Вас новое приложение✨",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontSize: 15),
                    ),
                    Text(
                      "Приносим извинения за доставленные неудобства🥲",
                      textAlign: TextAlign.center,
                      style: TextStyle(fontSize: 15),
                    )
                  ],
                ),
              ),
            ],
          ),
        )));
  }
}
