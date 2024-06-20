import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_rating_bar/flutter_rating_bar.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:mobile_web/helpers/api_caller.dart' as api_caller;
import 'package:mobile_web/helpers/render_helper.dart';
import 'package:mobile_web/helpers/session_helper.dart';
import 'package:mobile_web/qr_page.dart';

class FeedbackPage extends StatefulWidget {
  final api_caller.OrderAndFeedback currentOrderAndFeedback;
  late TextEditingController _feedBackTextController;
  late bool isReadonlyFeedback;

  FeedbackPage(this.currentOrderAndFeedback, {Key? key}) : super(key: key) {
    _feedBackTextController = TextEditingController(
        text: currentOrderAndFeedback.feedback?.feedbackText);

    isReadonlyFeedback = currentOrderAndFeedback.feedback != null;
  }

  @override
  State<FeedbackPage> createState() => _FeedbackPageState();
}

class _FeedbackPageState extends State<FeedbackPage> {
  final _form = GlobalKey<FormState>();

  double scoreQuality = 3;
  double scoreService = 3;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("Оставить отзыв"),
          backgroundColor: Colors.white,
          foregroundColor: Colors.black,
        ),
        body: SingleChildScrollView(
          child: Center(
            child: Column(
              children: <Widget>[
                Padding(
                  padding: const EdgeInsets.all(30),
                  child: Card(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(30),
                      ),
                      color: const Color.fromARGB(180, 62, 178, 178),
                      child: Padding(
                        padding: const EdgeInsets.all(20),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            RenderHelper(
                              check: widget.currentOrderAndFeedback.order,
                              needToRenderStar: false,
                            ),
                            const SizedBox(height: 10),
                            Form(
                              key: _form,
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text("Обслуживание",
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontWeight: FontWeight.bold,
                                          fontSize: 18)),
                                  RatingBar.builder(
                                    initialRating: widget
                                            .currentOrderAndFeedback
                                            .feedback
                                            ?.scoreService ??
                                        scoreService,
                                    minRating: 1,
                                    direction: Axis.horizontal,
                                    allowHalfRating: false,
                                    ignoreGestures: widget.isReadonlyFeedback,
                                    itemCount: 5,
                                    itemPadding: const EdgeInsets.symmetric(
                                        horizontal: 1.0),
                                    itemBuilder: (context, _) => const Icon(
                                      Icons.star,
                                      color: Colors.white,
                                    ),
                                    onRatingUpdate: (rating) {
                                      setState(() {
                                        scoreService = rating;
                                      });
                                    },
                                  ),
                                  const Text("Качество приготовления",
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontWeight: FontWeight.bold,
                                          fontSize: 18)),
                                  RatingBar.builder(
                                    initialRating: widget
                                            .currentOrderAndFeedback
                                            .feedback
                                            ?.scoreQuality ??
                                        scoreQuality,
                                    minRating: 1,
                                    direction: Axis.horizontal,
                                    allowHalfRating: false,
                                    itemCount: 5,
                                    ignoreGestures: widget.isReadonlyFeedback,
                                    itemPadding: const EdgeInsets.symmetric(
                                        horizontal: 1.0),
                                    itemBuilder: (context, _) => const Icon(
                                      Icons.star,
                                      color: Colors.white,
                                    ),
                                    onRatingUpdate: (rating) {
                                      setState(() {
                                        scoreQuality = rating;
                                      });
                                    },
                                  )
                                ],
                              ),
                            ),
                            const SizedBox(height: 10),
                            TextFormField(
                              minLines: 1,
                              maxLines: 5,
                              controller: widget._feedBackTextController,
                              keyboardType: TextInputType.multiline,
                              enabled: !widget.isReadonlyFeedback,
                              decoration: InputDecoration(
                                fillColor: Colors.white,
                                focusColor: Colors.white,
                                filled: true,
                                hintText: "Напишите развернутый отзыв\n\n\n",
                                hintStyle: const TextStyle(color: Colors.grey),
                                enabledBorder: OutlineInputBorder(
                                  borderSide: const BorderSide(
                                    width: 1,
                                    color: Color.fromARGB(255, 216, 216, 216),
                                  ),
                                  borderRadius: BorderRadius.circular(20.0),
                                ),
                              ),
                            ),
                          ],
                        ),
                      )),
                ),
                Container(
                    child: widget.currentOrderAndFeedback.feedback == null
                        ? ElevatedButton(
                            style: ElevatedButton.styleFrom(
                              foregroundColor: const Color(0xff3EB2B2),
                              backgroundColor: Colors.white,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(25),
                              ),
                              elevation: 15.0,
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 100, vertical: 20),
                            ),
                            onPressed: () async {
                              SmartDialog.showLoading(
                                  msg: "Отправляем ваш отзыв ✈");
                              // await api_caller.ApiCaller().addFeedback(
                              //     api_caller.Feedback(
                              //         id: -1,
                              //         checkid: widget.currentCheck.id!,
                              //         scorequality: scoreQuality.toInt(),
                              //         scoreservice: scoreService.toInt(),
                              //         feedback: widget._controller.text));

                              var clientId =
                                  await SessionHelper().getClientId();

                              var feedback = api_caller.FeedbackForOrder(
                                  widget.currentOrderAndFeedback.order.id,
                                  clientId,
                                  scoreQuality,
                                  scoreService,
                                  widget._feedBackTextController.text,
                                  jsonEncode(widget
                                      .currentOrderAndFeedback.order
                                      .toJson()));

                              var status = await api_caller.ApiCaller()
                                  .setFeedback(feedback);

                              if (status !=
                                  api_caller.TrySetFeedbackResStatus.Ok) {
                                SmartDialog.showToast(
                                    'не удалось сохранить отзыв(');
                              } else {
                                SmartDialog.showToast('Спасибо за отзыв 🫶');

                                await SessionHelper()
                                    .setFeedbackForOrder(feedback);
                              }

                              await SmartDialog.dismiss();
                              if (mounted) {
                                await Navigator.pushAndRemoveUntil(
                                    context,
                                    MaterialPageRoute(
                                        builder: (context) => const QrPage()),
                                    (route) => false);
                              }
                            },
                            child: const Text(
                              "Отправить отзыв",
                              style:
                                  TextStyle(fontSize: 16, color: Colors.black),
                            ),
                          )
                        : const Text("Отзыв успешно отправлен 👌",
                            style:
                                TextStyle(fontSize: 16, color: Colors.black))),
                const SizedBox(height: 10),
              ],
            ),
          ),
        ));
  }
}
