import 'package:vat_simple_frontend/models/product.model.dart';

class OrderItem {
  final Product product;
  int quantity;

  OrderItem({required this.product, this.quantity = 1});

  double get totalPrice => product.unitPrice * quantity;
}