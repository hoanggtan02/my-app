import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/product_service.api.dart';
import 'package:vat_simple_frontend/models/product.model.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/providers/cart_provider.dart';
import 'package:vat_simple_frontend/screens/invoices/checkout_screen.dart';

class ProductSelectionScreen extends StatefulWidget {
  const ProductSelectionScreen({super.key});

  @override
  State<ProductSelectionScreen> createState() => _ProductSelectionScreenState();
}

class _ProductSelectionScreenState extends State<ProductSelectionScreen> {
  late Future<List<Product>> _productsFuture;

  @override
  void initState() {
    super.initState();
    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    _productsFuture = ProductService().getProducts(token);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Chọn sản phẩm')),
      body: FutureBuilder<List<Product>>(
        future: _productsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError || !snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('Không thể tải sản phẩm.'));
          }
          final products = snapshot.data!;
          return GridView.builder(
            padding: const EdgeInsets.all(12),
            gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 200, // Chiều rộng tối đa của mỗi item
              childAspectRatio: 0.9,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
            ),
            itemCount: products.length,
            itemBuilder: (ctx, i) => ProductGridItem(product: products[i]),
          );
        },
      ),
      // Nút giỏ hàng
      floatingActionButton: Consumer<CartProvider>(
        builder: (_, cart, ch) => Badge(
          label: Text(cart.itemCount.toString()),
          isLabelVisible: cart.itemCount > 0,
          child: ch!,
        ),
        child: FloatingActionButton(
          onPressed: () {
            Navigator.of(context).push(
              MaterialPageRoute(builder: (ctx) => const CheckoutScreen()),
            );
          },
          child: const Icon(Icons.shopping_cart_outlined),
        ),
      ),
    );
  }
}

// trong file create_invoice_pos_screen.dart
class ProductGridItem extends StatelessWidget {
  const ProductGridItem({super.key, required this.product});
  final Product product;

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () {
          Provider.of<CartProvider>(context, listen: false).addItem(product);
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Đã thêm "${product.name}" vào giỏ hàng.'),
              duration: const Duration(seconds: 1),
            ),
          );
        },
        child: GridTile(
          footer: GridTileBar(
            backgroundColor: Colors.black54,
            title: Text(product.name, textAlign: TextAlign.center),
            subtitle: Text('${NumberFormat.decimalPattern().format(product.unitPrice)} đ', textAlign: TextAlign.center),
          ),
          // *** THAY THẾ ICON BẰNG HÌNH ẢNH ***
          child: (product.imageUrl != null && product.imageUrl!.isNotEmpty)
              ? Image.network(
                  product.imageUrl!,
                  fit: BoxFit.cover,
                  // Hiển thị loading trong khi tải ảnh
                  loadingBuilder: (context, child, loadingProgress) {
                    if (loadingProgress == null) return child;
                    return const Center(child: CircularProgressIndicator());
                  },
                  // Hiển thị icon nếu load ảnh lỗi
                  errorBuilder: (context, error, stackTrace) {
                    return Icon(Icons.image_not_supported, size: 60, color: Colors.grey[400]);
                  },
                )
              : Icon(Icons.fastfood, size: 60, color: Colors.grey[400]),
        ),
      ),
    );
  }
}