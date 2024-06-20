import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:shopping_cart/shopping_cart.dart';

import '../../../../helpers/api_caller.dart';
import '../../../cart/cart_item.dart';
import '../../product_detail/presentation/product_detail_page.dart';
import 'components/home_product_card_widget.dart';
import 'package:badges/badges.dart' as badges;

class ShopHomePage extends StatefulWidget {
  final String markerId;

  const ShopHomePage({super.key, required this.markerId});

  @override
  State<ShopHomePage> createState() => _ShopHomePageState();
}

class _ShopHomePageState extends State<ShopHomePage>
    with TickerProviderStateMixin {
  final Cart<CartItem> _cart = ShoppingCart.getInstance<CartItem>();
  late Future<List<ProductItem>> _futureProducts;
  String selectedCategory = "Главный экран";

  @override
  void initState() {
    super.initState();
    _futureProducts = ApiCaller().fetchProducts(widget.markerId);
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;
    return Scaffold(
      backgroundColor: const Color(0xFFF8F7F7),
      body: Column(
        children: [
          SizedBox(
            height: size.height * 0.11, // Adjust the height as needed
            child: Padding(
              padding: const EdgeInsets.only(top: 30, left: 16, right: 16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'ZebraCoffee Menu ${widget.markerId}',
                    style: GoogleFonts.poppins(
                      fontSize: 18, // Adjust the font size as needed
                      fontWeight: FontWeight.w600,
                      color: Colors.black,
                    ),
                  ),
                ],
              ),
            ),
          ),
          SizedBox(
            height: 50, // Adjust the height as needed
            child: FutureBuilder<List<ProductItem>>(
              future: _futureProducts,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                } else if (snapshot.hasError) {
                  return Center(child: Text('Error: ${snapshot.error}'));
                } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
                  return const Center(child: Text('No data available'));
                } else {
                  final categories = <String>{};

                  for (final product in snapshot.data!) {
                    if (product.category != null) {
                      categories.add(product.category!);
                    }
                  }

                  final sortedCategories = categories.toList()
                    ..sort((a, b) {
                      if (a == "Главный экран") {
                        return -1; // "Главный экран" comes first
                      } else if (b == "Главный экран") {
                        return 1; // "Главный экран" comes first
                      } else {
                        return a.compareTo(b);
                      }
                    });

                  return SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: Row(
                      children: sortedCategories.map((category) {
                        final isSelected = category == selectedCategory;
                        return Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 8),
                          child: ElevatedButton(
                            onPressed: () {
                              setState(() {
                                selectedCategory = category;
                              });
                              // Handle category selection
                            },
                            style: ButtonStyle(
                              backgroundColor: MaterialStateProperty.all<Color>(
                                isSelected
                                    ? const Color.fromARGB(255, 62, 178, 178)
                                    : Colors.white,
                              ),
                            ),
                            child: Text(
                              category,
                              style: TextStyle(
                                color: isSelected ? Colors.white : Colors.black,
                              ),
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  );
                }
              },
            ),
          ),
          const SizedBox(
            height: 10,
          ),
          Expanded(
            child: Stack(
              children: [
                Container(
                  clipBehavior: Clip.antiAlias,
                  decoration: const BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(40),
                      topRight: Radius.circular(40),
                    ),
                  ),
                  child: FutureBuilder<List<ProductItem>>(
                    future: _futureProducts,
                    builder: (context, snapshot) {
                      if (snapshot.connectionState == ConnectionState.waiting) {
                        return const Center(child: CircularProgressIndicator());
                      } else if (snapshot.hasError) {
                        return Center(child: Text('Error: ${snapshot.error}'));
                      } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
                        return const Center(child: Text('No data available'));
                      } else {
                        final products = snapshot.data!;

                        // If "Главный экран" category is selected, show all products
                        final selectedProducts =
                            selectedCategory == "Главный экран"
                                ? products
                                : products
                                    .where((product) =>
                                        product.category == selectedCategory)
                                    .toList();

                        return ListView.separated(
                          itemCount: selectedProducts.length,
                          separatorBuilder: (context, index) =>
                              const SizedBox(height: 10),
                          itemBuilder: (context, index) {
                            final product = selectedProducts[index];
                            return GestureDetector(
                              onTap: () async {
                                await Navigator.push(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) => ShopProductDetailPage(
                                      image: 'assets/4.png',
                                      nabors: product.nabors,
                                      name: product.itemName,
                                      price: product.price,
                                      productID: product.itemId,
                                    ),
                                  ),
                                ).then((value) => setState(() {}));
                              },
                              child: HomeProductCardWidget(
                                name: product.itemName,
                                price: product.price.roundToDouble(),
                                image: 'assets/4.png', // Modify as needed
                              ),
                            );
                          },
                        );
                      }
                    },
                  ),
                ),
                _cart.itemCount == 0
                    ? Container()
                    : Align(
                        alignment: Alignment.bottomCenter,
                        child: Container(
                          width: double.infinity,
                          height: 72,
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.8),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.grey.withOpacity(0.1),
                                spreadRadius: 2,
                                blurRadius: 20,
                                offset: const Offset(0, -10),
                              ),
                            ],
                          ),
                          child: ClipRRect(
                            child: BackdropFilter(
                              filter: ImageFilter.blur(sigmaX: 8, sigmaY: 8),
                              child: Container(
                                decoration: BoxDecoration(
                                  color: Colors.white.withOpacity(0.2),
                                  borderRadius: BorderRadius.circular(20),
                                ),
                                child: Stack(
                                  children: [
                                    Center(
                                      child: Row(
                                        mainAxisAlignment:
                                            MainAxisAlignment.spaceAround,
                                        children: [
                                          Text(
                                            "${_cart.cartTotal} тг",
                                            style: const TextStyle(
                                                fontWeight: FontWeight.bold,
                                                fontSize: 18),
                                          ),
                                          SizedBox(
                                            width: 200,
                                            height: 50,
                                            child: ElevatedButton(
                                              style: ElevatedButton.styleFrom(
                                                foregroundColor: Colors.white,
                                                backgroundColor:
                                                    const Color.fromARGB(
                                                        255, 62, 178, 178),
                                                shape: RoundedRectangleBorder(
                                                  borderRadius:
                                                      BorderRadius.circular(25),
                                                ),
                                                elevation: 0,
                                              ),
                                              onPressed: () async {
                                                showModalBottomSheet<void>(
                                                  context: context,
                                                  builder:
                                                      (BuildContext context) {
                                                    return getCartModal(
                                                        context);
                                                  },
                                                );
                                              },
                                              child: Row(
                                                mainAxisAlignment:
                                                    MainAxisAlignment.center,
                                                children: [
                                                  const Text(
                                                    "оплатить",
                                                    style:
                                                        TextStyle(fontSize: 16),
                                                  ),
                                                  badges.Badge(
                                                    badgeContent: Text(
                                                      _cart.itemCount
                                                          .toString(),
                                                      style: const TextStyle(
                                                          color: Colors.white),
                                                    ),
                                                    child: const FaIcon(
                                                      FontAwesomeIcons
                                                          .cartShopping,
                                                      color: Colors.white,
                                                    ),
                                                  ),
                                                ],
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ),
                        ),
                      )
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget getCartModal(BuildContext context) {
    return SizedBox(
      height: 300,
      child: Column(
        children: <Widget>[
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              IconButton(
                icon: const Icon(Icons.close),
                onPressed: () => Navigator.pop(context),
              ),
            ],
          ),
          Column(
              children: _cart.cartItems
                  .map((e) => Padding(
                        padding: const EdgeInsets.all(5.0),
                        child:
                            Text("${e.name} x ${e.quantity} - ${e.price} тг"),
                      ))
                  .toList()),
          const Spacer(),
          SizedBox(
            width: 200,
            height: 50,
            child: ElevatedButton(
              style: ElevatedButton.styleFrom(
                foregroundColor: Colors.white,
                backgroundColor: Colors.red,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(25),
                ),
                elevation: 0,
              ),
              onPressed: () async {},
              child: const Text(
                "оплатить c kaspi.kz",
                style: TextStyle(fontSize: 16),
              ),
            ),
          ),
          TextButton(
              onPressed: () {
                _cart.clearCart();
                setState(() {});
                Navigator.pop(context);
              },
              child: const Text("очистить корзину"))
        ],
      ),
    );
  }
}
