import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:mobile_web/helpers/api_caller.dart';

class RenderHelper extends StatefulWidget {
  final Order check;
  final bool needToRenderStar;

  RenderHelper({required this.check, required this.needToRenderStar});

  @override
  _RenderHelperState createState() => _RenderHelperState();
}

class _RenderHelperState extends State<RenderHelper> {
  bool showMoreItems = false;

  @override
  Widget build(BuildContext context) {
    List<Widget> itemsInOrder = <Widget>[];
    var title = Row(
      children: [
        Text(
          DateFormat('d MMMM, yyyy HH:mm', 'ru').format(widget.check.openedat),
          style: const TextStyle(fontSize: 16, color: Colors.white),
        ),
        const Spacer(),
        Text(
          "#${widget.check.id}",
          style: const TextStyle(fontSize: 16, color: Colors.white),
        ),
      ],
    );

    for (var x in widget.check.tovars) {
      var text = "- ${x.name}";
      if (x.quantity > 1) {
        text += " x ${x.quantity}";
      }

      itemsInOrder.add(Align(
          alignment: Alignment.topLeft,
          child: FittedBox(
            child: Text(
              text,
              overflow: TextOverflow.ellipsis,
              softWrap: false,
              style: const TextStyle(color: Colors.white, fontSize: 18),
            ),
          )));
    }

    for (var x in widget.check.techCarts) {
      var text = " ${x.name}";
      if (x.quantity > 1) {
        text += " x ${x.quantity}";
      }

      itemsInOrder.add(Align(
          alignment: Alignment.topLeft,
          child: Text(
            text,
            style: const TextStyle(
                color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
          )));

      if (x.modificators.isNotEmpty) {
        var modificatorsText = "\t\t (${x.modificators.join(", ")})";

        itemsInOrder.add(Align(
            alignment: Alignment.topLeft,
            child: Text(
              modificatorsText,
              style: const TextStyle(color: Colors.white, fontSize: 12),
            )));
      }
    }
    var maxItemsToShow = showMoreItems ? itemsInOrder.length : 2;
    var visibleItems = itemsInOrder.sublist(0, maxItemsToShow);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(right: 10, bottom: 10),
          child: title,
        ),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: visibleItems,
        ),
        if (itemsInOrder.length >
            2) // Show "Show More" button when there are more than 2 items
          Row(
            children: [
              TextButton(
                onPressed: () {
                  setState(() {
                    showMoreItems = !showMoreItems;
                  });
                },
                child: const Text(
                  "...",
                  style: TextStyle(
                      color: Colors.white,
                      fontSize: 36,
                      fontWeight: FontWeight.bold),
                ),
              ),
              const Spacer(),
              widget.needToRenderStar
                  ? const Padding(
                      padding: EdgeInsets.all(8.0),
                      child: Icon(
                        Icons.star,
                        color: Colors.white,
                      ),
                    )
                  : Container()
            ],
          ),
      ],
    );
  }
}
