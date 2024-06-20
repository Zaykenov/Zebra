import 'package:mobile_web/shoping/config/theme/app_colors.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import '../../../../shared/widgets/app_button.dart';

class HomeProductCardWidget extends StatelessWidget {
  final String image;
  final String name;
  final double price;
  const HomeProductCardWidget({
    super.key,
    required this.image,
    required this.name,
    required this.price,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
      child: Stack(
        alignment: Alignment.center,
        children: [
          Column(
            children: [
              Row(
                children: [
                  const SizedBox(
                    width: 120,
                  ),
                  Expanded(
                    child: SizedBox(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            name,
                            style: GoogleFonts.poppins(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                              color: AppColors.secondaryColor,
                            ),
                          ),
                          const SizedBox(height: 30),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
              Container(
                height: 68,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.circular(32),
                    bottomRight: Radius.circular(32),
                    topRight: Radius.circular(6),
                    bottomLeft: Radius.circular(6),
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(.06),
                      blurRadius: 8,
                      offset: const Offset(0, 10),
                    ),
                  ],
                ),
                child: Row(
                  children: [
                    const SizedBox(
                      width: 120,
                    ),
                    SizedBox(
                      width: 86,
                      child: AppButton(
                        onPressed: () {},
                        reverse: true,
                        child: Text(
                          '$price тг',
                          style: GoogleFonts.poppins(
                              fontSize: 14,
                              fontWeight: FontWeight.bold,
                              color: Colors.white),
                        ),
                      ),
                    ),
                    const Spacer(),
                    // Padding(
                    //   padding: const EdgeInsets.only(right: 30),
                    //   child: Row(
                    //     children: [
                    //       FaIcon(
                    //         FontAwesomeIcons.solidStar,
                    //         color: Colors.yellow[700],
                    //         size: 16,
                    //       ),
                    //       const SizedBox(width: 4),
                    //       Text(
                    //         '4.5',
                    //         style: GoogleFonts.poppins(
                    //           fontSize: 16,
                    //           fontWeight: FontWeight.bold,
                    //           color: Colors.yellow[700],
                    //         ),
                    //       ),
                    //     ],
                    //   ),
                    // ),
                    //const SizedBox(width: 16),
                  ],
                ),
              )
            ],
          ),
          Positioned(
            left: -0,
            child: SizedBox(
              height: 140,
              width: 120,
              child: Hero(
                tag: image,
                child: Image.asset(
                  image,
                  fit: BoxFit.fitHeight,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
