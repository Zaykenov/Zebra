import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:mobile_web/main.dart';
import 'package:mobile_web/support_service_page.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:url_launcher/url_launcher.dart';

class AboutPage extends StatefulWidget {
  const AboutPage({Key? key}) : super(key: key);
  @override
  State<AboutPage> createState() => _AboutPageState();
}

class _AboutPageState extends State<AboutPage> {
  String appVersion = "";

  @override
  void initState() {
    super.initState();
    startLoadAppVersion();
  }

  Future startLoadAppVersion() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    String version = packageInfo.version;
    setState(() {
      appVersion = version;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("О программе"),
          backgroundColor: Colors.white,
          foregroundColor: Colors.black,
        ),
        body: SingleChildScrollView(
          child: Center(
            child: Column(
              children: <Widget>[
                InkWell(
                  onTap: () async {
                    final uri = Uri.parse('https://zebra-crm.kz/privacy');
                    await launchUrl(uri, mode: LaunchMode.externalApplication);
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
                              Icons.lock,
                            ),
                            SizedBox(width: 5),
                            Text("Политика конфиденциальности")
                          ],
                        ),
                        //child:,
                      ),
                    ),
                  ),
                ),
                //const SizedBox(height: 10),
                InkWell(
                  onTap: () async {
                    final uri = Uri.parse('https://zebra-crm.kz/terms');
                    await launchUrl(uri, mode: LaunchMode.externalApplication);
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
                              Icons.handshake,
                            ),
                            SizedBox(width: 5),
                            Text("Пользовательское соглашение")
                          ],
                        ),
                      ),
                    ),
                  ),
                ),

                InkWell(
                  onTap: () async {
                    await Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => const SupportServicePage()));
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
                              Icons.support,
                            ),
                            SizedBox(width: 5),
                            Text("Служба поддержки")
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Text("версия приложения: "),
                    GestureDetector(
                        onLongPress: () async {
                          var appConfig = await SessionHelper().getAppConfig();

                          AppConfig newConfig;
                          if (appConfig.envName == EnvName.production) {
                            newConfig = AppConfig(
                                apiUrl: ApiCaller().stagingApiUrl,
                                envName: EnvName.staging);
                          } else {
                            newConfig = AppConfig(
                                apiUrl: ApiCaller().prodApiUrl,
                                envName: EnvName.production);
                          }

                          await SessionHelper().setAppConfig(newConfig);
                          SmartDialog.showToast("${newConfig.envName}");
                          if (mounted) {
                            await SessionHelper().userSignOut().then(
                                (value) => RestartWidget.restartApp(context));
                          }
                        },
                        child: Text(appVersion)),
                  ],
                )
              ],
            ),
          ),
        ));
  }
}
