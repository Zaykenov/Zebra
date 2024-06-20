import 'package:flutter/material.dart';
import 'package:mobile_web/helpers/api_caller.dart';

class ModifierCounter extends StatefulWidget {
  final Modificator modificator;
  final double price;
  final Function() onModifierStateChanged;

  const ModifierCounter(
      {super.key,
      required this.modificator,
      required this.price,
      required this.onModifierStateChanged});

  @override
  _ModifierCounterState createState() => _ModifierCounterState();
}

class _ModifierCounterState extends State<ModifierCounter> {
  void _incrementPressCount() {
    setState(() {
      widget.modificator.isPressed = true; // Toggle the flag on each press
      if (widget.modificator.isPressed) {
        widget.modificator.taps++;
      }
    });
    widget.onModifierStateChanged();
  }

  void _closeModification() {
    setState(() {
      widget.modificator.isPressed = false; // Toggle the flag on each press
      widget.modificator.taps = 0;
    });
    widget.onModifierStateChanged();
  }

  @override
  Widget build(BuildContext context) {
    final containerColor =
        widget.modificator.isPressed ? Colors.black : Colors.white;
    final textColor =
        widget.modificator.isPressed ? Colors.white : Colors.black;

    return GestureDetector(
      onTap: _incrementPressCount,
      child: Container(
        padding: const EdgeInsets.all(16.0),
        decoration: BoxDecoration(
          color: containerColor,
          borderRadius: BorderRadius.circular(30.0),
          boxShadow: [
            BoxShadow(
              color: Colors.grey.withOpacity(0.3),
              spreadRadius: 2,
              blurRadius: 4,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                GestureDetector(
                  onTap: _closeModification,
                  child: Visibility(
                    visible: widget.modificator.isPressed,
                    child: Container(
                      width: 30,
                      height: 30,
                      decoration: const BoxDecoration(
                        color: Colors.white,
                        shape: BoxShape.circle,
                      ),
                      child: const Center(
                        child: Icon(Icons.close, color: Colors.red),
                      ),
                    ),
                  ),
                ),
                Visibility(
                  visible: widget.modificator.isPressed,
                  child: Container(
                    width: 30,
                    height: 30,
                    decoration: const BoxDecoration(
                      color: Colors.red,
                      shape: BoxShape.circle,
                    ),
                    child: Center(
                      child: Text(
                        '${widget.modificator.taps}',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(
              height: 20,
            ),
            Expanded(
              child: Text(
                widget.modificator.ingredientName,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: textColor,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Container(
              margin: const EdgeInsets.symmetric(horizontal: 8),
              padding: const EdgeInsets.symmetric(
                vertical: 4,
                horizontal: 12,
              ),
              decoration: BoxDecoration(
                color: const Color(0xff3EB2B2),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Text(
                '${widget.price.round()} тг',
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
