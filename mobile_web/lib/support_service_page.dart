import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:mobile_web/main.dart';
import 'package:url_launcher/url_launcher.dart';

class SupportServicePage extends StatefulWidget {
  const SupportServicePage({Key? key}) : super(key: key);
  @override
  State<SupportServicePage> createState() => _SupportServicePageState();
}

class _SupportServicePageState extends State<SupportServicePage> {
  final _supportPhoneNumber = "+7 771 458 29 18";
  final _supportMail = "saparbek394@gmail.com";

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("Служба поддержки"),
          backgroundColor: Colors.white,
          foregroundColor: Colors.black,
        ),
        body: SingleChildScrollView(
          child: Center(
            child: Column(
              children: <Widget>[
                SizedBox(
                  width: double.infinity,
                  child: Card(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(2),
                    ),
                    elevation: 0,
                    child: const Padding(
                      padding: EdgeInsets.all(20.0),
                      child: Column(
                        children: [
                          Image(image: AssetImage('assets/logo.png')),
                          Text(
                            "Мы дорожим нашими клиентами и всегда готовы помочь",
                            textAlign: TextAlign.center,
                            style: TextStyle(fontSize: 16),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                //const SizedBox(height: 10),
                // InkWell(
                //   onTap: () async {
                //     await launchUrl(Uri.parse("tel://$_supportPhoneNumber"));
                //   },
                //   child: SizedBox(
                //     width: double.infinity,
                //     child: Card(
                //       shape: RoundedRectangleBorder(
                //         borderRadius: BorderRadius.circular(2),
                //       ),
                //       child: Padding(
                //         padding: const EdgeInsets.all(20.0),
                //         child: Row(
                //           children: [
                //             const Icon(
                //               Icons.phone,
                //             ),
                //             const SizedBox(width: 5),
                //             Text(_supportPhoneNumber),
                //           ],
                //         ),
                //       ),
                //     ),
                //   ),
                // ),

                InkWell(
                  onTap: () async {
                    await launchUrl(Uri.parse("mailto:$_supportPhoneNumber"));
                  },
                  child: SizedBox(
                    width: double.infinity,
                    child: Card(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(2),
                      ),
                      child: Padding(
                        padding: const EdgeInsets.all(20.0),
                        child: Row(
                          children: [
                            const Icon(
                              Icons.mail,
                            ),
                            const SizedBox(width: 5),
                            Text(_supportMail),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),

                InkWell(
                  onTap: () async {
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
                            const Padding(
                              padding: EdgeInsets.only(top: 10),
                              child: Text(
                                  'Вы действительно хотите удалить свой аккаунт безвозвратно?',
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      fontWeight: FontWeight.bold,
                                      fontSize: 18)),
                            ),
                            const SizedBox(
                              height: 10,
                            ),
                            TextButton(
                              style: ButtonStyle(
                                foregroundColor:
                                    MaterialStateProperty.all<Color>(
                                        Colors.blue),
                              ),
                              onPressed: () async {
                                var clientId =
                                    await SessionHelper().getClientId();
                                var removeStatus =
                                    await ApiCaller().tryRemoveUser(clientId);
                                if (removeStatus == TryRemoveUserResStatus.Ok) {
                                  await SessionHelper().userSignOut().then(
                                      (value) =>
                                          RestartWidget.restartApp(context));
                                } else {
                                  SmartDialog.showToast(
                                      "не удалось удалить аккаунт");
                                }
                              },
                              child: const Text('Удалить мой аккаунт',
                                  style: TextStyle(color: Colors.red)),
                            )
                          ],
                        ),
                      );
                    });
                  },
                  child: SizedBox(
                    width: double.infinity,
                    child: Card(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(2),
                      ),
                      child: const Padding(
                        padding: EdgeInsets.all(20.0),
                        child: Row(
                          children: [
                            Icon(
                              Icons.person_remove,
                            ),
                            SizedBox(width: 5),
                            Text("Удалить мой аккаунт"),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ));
  }
}
