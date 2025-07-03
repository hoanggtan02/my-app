import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/cart_provider.dart';
import 'package:vat_simple_frontend/api/invoice_service.api.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/models/order_item.model.dart';

class CheckoutScreen extends StatelessWidget {
  const CheckoutScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final cart = Provider.of<CartProvider>(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Hóa đơn & Thanh toán')),
      body: Column(
        children: [
          Expanded(
            child: cart.items.isEmpty
                ? const Center(child: Text('Giỏ hàng của bạn trống.'))
                : ListView.builder(
                    padding: const EdgeInsets.all(8),
                    itemCount: cart.items.length,
                    itemBuilder: (ctx, i) => CheckoutItemCard(item: cart.items[i]),
                  ),
          ),
          // Phần tổng tiền và nút thanh toán
          Card(
            margin: const EdgeInsets.all(15),
            child: Padding(
              padding: const EdgeInsets.all(8),
              child: Column(
                children: [
                  _buildTotalRow('Trước thuế:', '${NumberFormat.decimalPattern().format(cart.subtotal)} VND'),
                  _buildTotalRow('Thuế (10%):', '${NumberFormat.decimalPattern().format(cart.tax)} VND'),
                  const Divider(),
                  _buildTotalRow('TỔNG CỘNG:', '${NumberFormat.decimalPattern().format(cart.total)} VND', isTotal: true),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                onPressed: cart.items.isEmpty
                    ? null
                    : () async {
                        final token = Provider.of<AuthProvider>(context, listen: false).token!;
                        try {
                          await InvoiceService().createInvoice(token, cart.items);

                          // Xóa giỏ hàng và quay về màn hình chọn sản phẩm
                          cart.clearCart();
                          if (context.mounted) {
                            Navigator.of(context).pop();
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('Thanh toán thành công!'), backgroundColor: Colors.green),
                            );
                          }
                        } catch (e) {
                          if (context.mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text(e.toString())),
                            );
                          }
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 12),
                        backgroundColor: Theme.of(context).primaryColor,
                        foregroundColor: Colors.white,
                      ),
                      child: const Text('XÁC NHẬN THANH TOÁN', style: TextStyle(fontSize: 16)),
                    ),
                  ),
                ],
              ),
            ),
          )
        ],
      ),
    );
  }

Widget _buildTotalRow(String title, String value, {bool isTotal = false}) {
  final style = TextStyle(
    fontSize: isTotal ? 20 : 16,
    fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
  );
  return Padding(
    padding: const EdgeInsets.symmetric(vertical: 2.0),
    child: Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(title, style: style),
        Text(value, style: style)
      ],
    ),
  );
}
}

class CheckoutItemCard extends StatelessWidget {
  const CheckoutItemCard({super.key, required this.item});
  final OrderItem item;

  @override
  Widget build(BuildContext context) {
    final cart = Provider.of<CartProvider>(context, listen: false);
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: ListTile(
        leading: CircleAvatar(child: Icon(Icons.fastfood_outlined)),
        title: Text(item.product.name),
        subtitle: Text('${NumberFormat.decimalPattern().format(item.product.unitPrice)} đ'),
        trailing: SizedBox(
          width: 120,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              IconButton(icon: Icon(Icons.remove), onPressed: () => cart.updateQuantity(item.product.id, item.quantity - 1)),
              Text(item.quantity.toString(), style: TextStyle(fontSize: 16)),
              IconButton(icon: Icon(Icons.add), onPressed: () => cart.updateQuantity(item.product.id, item.quantity + 1)),
            ],
          ),
        ),
      ),
    );
  }
}