import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/providers/cart_provider.dart';
import 'package:vat_simple_frontend/screens/invoices/product_selection_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthProvider>(context, listen: false);

    return Scaffold(
      appBar: AppBar(
        title: Text('Trang chủ - ${authProvider.user?.companyName ?? ''}'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            tooltip: 'Đăng xuất',
            onPressed: () {
              Provider.of<AuthProvider>(context, listen: false).logout();
            },
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Chào mừng bạn đến với trang quản lý!', style: TextStyle(fontSize: 22)),
            const SizedBox(height: 30),
            ElevatedButton.icon(
              icon: const Icon(Icons.point_of_sale),
              label: const Text('Mở Giao diện Bán hàng'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 15),
                textStyle: const TextStyle(fontSize: 18),
              ),
              onPressed: () {
                Provider.of<CartProvider>(context, listen: false).clearCart();
                Navigator.of(context).push(MaterialPageRoute(
                  builder: (context) => const ProductSelectionScreen(),
                ));
              },
            ),
          ],
        ),
      ),
    );
  }
}