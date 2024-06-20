import 'package:shopping_cart/shopping_cart.dart';

class CartItem extends ItemModel {
  final String name;

  CartItem(this.name,
      {required super.id, required super.price, required super.quantity});

  @override
  bool operator ==(Object other) =>
      other is CartItem && other.runtimeType == runtimeType && other.id == id;

  @override
  int get hashCode => id.hashCode;
}
