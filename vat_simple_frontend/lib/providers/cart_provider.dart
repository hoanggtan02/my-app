import 'package:flutter/material.dart';
import 'package:vat_simple_frontend/models/order_item.model.dart';
import 'package:vat_simple_frontend/models/product.model.dart';

// Provider quản lý toàn bộ giỏ hàng
class CartProvider with ChangeNotifier {
  final List<OrderItem> _items = [];

  List<OrderItem> get items => [..._items];
  int get itemCount => _items.length;

  double get subtotal => _items.fold(0, (sum, item) => sum + item.totalPrice);
  double get tax => subtotal * 0.10;
  double get total => subtotal + tax;

  void addItem(Product product) {
    final existingItemIndex = _items.indexWhere((item) => item.product.id == product.id);
    if (existingItemIndex >= 0) {
      _items[existingItemIndex].quantity++;
    } else {
      _items.add(OrderItem(product: product));
    }
    notifyListeners();
  }

  void removeItem(String productId) {
    _items.removeWhere((item) => item.product.id == productId);
    notifyListeners();
  }
  
  void updateQuantity(String productId, int newQuantity) {
     final existingItemIndex = _items.indexWhere((item) => item.product.id == productId);
     if (existingItemIndex >= 0) {
       if (newQuantity > 0) {
         _items[existingItemIndex].quantity = newQuantity;
       } else {
         removeItem(productId);
       }
       notifyListeners();
     }
  }

  void clearCart() {
    _items.clear();
    notifyListeners();
  }
}