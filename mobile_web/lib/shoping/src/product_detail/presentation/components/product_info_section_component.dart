import 'package:flutter/material.dart';
import 'package:mobile_web/helpers/api_caller.dart';

import 'line_product_info_widget.dart';

class ProductInfoSectionComponent extends StatelessWidget {
  final Nabor selectedNabor;
  const ProductInfoSectionComponent({super.key, required this.selectedNabor});

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.only(left: 52),
      child: Column(
        children: [
          LineProductInfoWidget(
            title: 'HOPS',
            description: '2 Row, Torrified Wheat',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'MALTS',
            description: 'Cascade, First Gold, Mt.Hood',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'HOPS',
            description: '2 Row, Torrified Wheat',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'MALTS',
            description: 'Cascade, First Gold, Mt.Hood',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'HOPS',
            description: '2 Row, Torrified Wheat',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'MALTS',
            description: 'Cascade, First Gold, Mt.Hood',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'MALTS',
            description: 'Cascade, First Gold, Mt.Hood',
          ),
          SizedBox(height: 14),
          LineProductInfoWidget(
            title: 'MALTS',
            description: 'Cascade, First Gold, Mt.Hood',
          ),
          SizedBox(height: 14),
        ],
      ),
    );
  }
}
