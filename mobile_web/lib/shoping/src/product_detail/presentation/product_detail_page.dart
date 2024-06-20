import 'package:mobile_web/helpers/api_caller.dart';
import 'package:mobile_web/shoping/config/theme/app_colors.dart';
import 'package:mobile_web/shoping/shared/widgets/app_bar.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:mobile_web/shoping/src/product_detail/presentation/components/modificator_counter.dart';
import 'package:shopping_cart/shopping_cart.dart';

import '../../../cart/cart_item.dart';
import '../../../shared/widgets/app_button.dart';
import 'components/product_counter_widget.dart';

class ShopProductDetailPage extends StatefulWidget {
  final String image;
  final List<Nabor>? nabors;
  final String name;
  final double price;
  final int productID;
  const ShopProductDetailPage({
    super.key,
    required this.image,
    required this.nabors,
    required this.name,
    required this.price,
    required this.productID,
  });

  @override
  State<ShopProductDetailPage> createState() => _ShopProductDetailPageState();
}

class _ShopProductDetailPageState extends State<ShopProductDetailPage> {
  final DraggableScrollableController _scrollableController =
      DraggableScrollableController();
  final ValueNotifier<double> _circleRadius = ValueNotifier<double>(180);
  final ValueNotifier<double> _imageHeight = ValueNotifier<double>(250);
  final double _minSize = 0.44;
  final double _maxSize = 0.76;

  final ValueNotifier<bool> _buttonLoading = ValueNotifier(false);
  late Nabor selectedNabor;
  int _selectedQuantity = 1;
  double modificatorsPrice = 0;

  double _calculateTotalPrice() {
    double totalPrice =
        (widget.price + _modificatorsPrice()) * _selectedQuantity;

    return totalPrice;
  }

  double _modificatorsPrice() {
    double modificatorsPrice = 0;

    if (widget.nabors != null && widget.nabors!.isNotEmpty) {
      for (Nabor nabors in widget.nabors!) {
        for (Modificator modifier in nabors.modificators) {
          if (modifier.isPressed) {
            modificatorsPrice += nabors.price.round() * modifier.taps;
          }
        }
      }
    }

    return modificatorsPrice;
  }

  void _handleModifierStateChanged() {
    setState(() {
      // Update the state related to the selected modifiers here
      // For example, you can update the `isPressed` and `taps` properties of modifiers

      // Recalculate the modificator prices

      // You can use this calculated price as needed or update any other state variables
    });
  }

  void _resetModifiers() {
    for (Nabor nabor in widget.nabors ?? []) {
      for (Modificator modifier in nabor.modificators) {
        modifier.isPressed = false;
        modifier.taps = 0;
      }
    }
  }

  @override
  void initState() {
    super.initState();
    _scrollableController.addListener(_updateCircleRadius);
    if (widget.nabors != null && widget.nabors!.isNotEmpty) {
      selectedNabor = widget.nabors!.first;
    } else {
      // Initialize with a default value when nabors list is empty
      selectedNabor = Nabor(
          name: 'Default Nabor',
          modificators: [],
          description: '',
          max: 0,
          min: 0,
          naborId: 0,
          price: 0);
    }
  }

  @override
  void dispose() {
    _scrollableController.removeListener(_updateCircleRadius);
    _resetModifiers();
    super.dispose();
  }

  void _updateCircleRadius() {
    double size = _scrollableController.size.clamp(_minSize, _maxSize);
    double normalizedSize = (size - _minSize) / (_maxSize - _minSize);
    _circleRadius.value = 180 - (normalizedSize * 100);
    _imageHeight.value = 250 - (normalizedSize * 100);
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 212, 212, 212),
      appBar: AppBarWidget(
        onBackPressed: () => Navigator.of(context).pop(),
      ),
      body: CustomScrollView(
        slivers: [
          SliverToBoxAdapter(
            child: Column(
              children: [
                const SizedBox(
                  height: 30,
                ),
                ValueListenableBuilder(
                  valueListenable: _imageHeight,
                  builder: (context, value, _) {
                    return Center(
                      child: SizedBox(
                        height: _imageHeight.value,
                        child: Hero(
                          tag: widget.image,
                          child: Image.asset(
                            widget.image,
                            fit: BoxFit.fitHeight,
                          ),
                        ),
                      ),
                    );
                  },
                ),
              ],
            ),
          ),
          SliverList(
            delegate: SliverChildListDelegate([
              const SizedBox(height: 15),
              Center(
                child: Text(
                  widget.name,
                  style: GoogleFonts.poppins(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: AppColors.secondaryColor,
                  ),
                ),
              ),
              const SizedBox(
                height: 10,
              ),
              Center(
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    'Кофе раф (или латте раф) – это изысканный напиток, который объединяет в себе богатство ароматного кофе и нежную кремовую текстуру молока. ',
                    style: GoogleFonts.poppins(
                      fontSize: 12,
                      fontWeight: FontWeight.normal,
                      color: AppColors.secondaryColor,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),
              ),
              const SizedBox(
                height: 10,
              ),
              if (widget.nabors != null && widget.nabors!.isNotEmpty)
                SizedBox(
                  height: 50,
                  child: ListView.builder(
                    scrollDirection: Axis.horizontal,
                    itemCount: widget.nabors!.length,
                    itemBuilder: (context, index) {
                      // ... (your ListView.builder code)
                      final nabor = widget.nabors![index];
                      return GestureDetector(
                        onTap: () {
                          setState(() {
                            selectedNabor = nabor;
                          });
                        },
                        child: Container(
                          margin: const EdgeInsets.symmetric(horizontal: 8),
                          padding: const EdgeInsets.symmetric(
                            vertical: 4,
                            horizontal: 12,
                          ),
                          decoration: BoxDecoration(
                            color: (selectedNabor == nabor)
                                ? AppColors.primaryColor
                                : Colors.white,
                            borderRadius: BorderRadius.circular(20),
                          ),
                          child: Center(
                            child: Text(
                              nabor.name,
                              style: TextStyle(
                                color: (selectedNabor == nabor)
                                    ? Colors.white
                                    : Colors.black,
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                ),
              const SizedBox(
                height: 10,
              ),
              SizedBox(
                height: 400, // Adjust the height as needed
                child: GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2, // Number of columns in the grid
                    mainAxisSpacing: 15.0, // Spacing between rows
                    crossAxisSpacing: 20.0, // Spacing between columns
                  ),
                  itemCount: selectedNabor.modificators.length,
                  itemBuilder: (context, index) {
                    final modifier = selectedNabor.modificators[index];
                    return ModifierCounter(
                      modificator: modifier,
                      price: selectedNabor.price,
                      onModifierStateChanged:
                          _handleModifierStateChanged, // Add this line
                    );
                  },
                ),
              ),
              // ... (other SliverChildListDelegate widgets if needed)
            ]),
          ),
        ],
      ),
      bottomNavigationBar: Container(
        color: Colors.white,
        padding: const EdgeInsets.all(16),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            ProductCounter(
              initialValue: _selectedQuantity,
              onValueChanged: (value) {
                setState(() {
                  _selectedQuantity = value;
                });
              },
            ),
            SizedBox(
              height: 60,
              width: size.width * .4,
              child: ValueListenableBuilder(
                valueListenable: _buttonLoading,
                builder: (_, isLoading, __) {
                  return AppButton(
                    reverse: true,
                    loading: isLoading,
                    onPressed: () {
                      if (isLoading) {
                        _buttonLoading.value = false;
                        return;
                      }
                      _buttonLoading.value = true;

                      final cart = ShoppingCart.getInstance<CartItem>();
                      cart.addItemToCart(CartItem(
                        widget.name,
                        id: widget.productID,
                        price: _calculateTotalPrice(),
                        quantity: _selectedQuantity,
                      ));
                      Navigator.of(context).pop();
                    },
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        Text(
                          '${_calculateTotalPrice().roundToDouble()} тг',
                          style: GoogleFonts.poppins(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                        Container(
                          width: 3,
                          height: 3,
                          decoration: const BoxDecoration(
                            color: Colors.white,
                            shape: BoxShape.circle,
                          ),
                        ),
                        Text(
                          'Добавить',
                          style: GoogleFonts.poppins(
                            fontSize: 12,
                            color: Colors.white,
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
