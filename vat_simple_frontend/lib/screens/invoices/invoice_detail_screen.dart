import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/invoice_service.api.dart';
import 'package:vat_simple_frontend/models/invoice.model.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';

class InvoiceDetailScreen extends StatefulWidget {
  final String invoiceId;
  const InvoiceDetailScreen({super.key, required this.invoiceId});

  @override
  State<InvoiceDetailScreen> createState() => _InvoiceDetailScreenState();
}

class _InvoiceDetailScreenState extends State<InvoiceDetailScreen> {
  late Future<Invoice> _invoiceDetailFuture;

  @override
  void initState() {
    super.initState();
    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    _invoiceDetailFuture = InvoiceService().getInvoiceDetails(token, widget.invoiceId);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Chi tiết Hóa đơn'),
      ),
      body: FutureBuilder<Invoice>(
        future: _invoiceDetailFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError || !snapshot.hasData) {
            return Center(child: Text('Lỗi: ${snapshot.error ?? 'Không có dữ liệu'}'));
          }

          final invoice = snapshot.data!;
          return ListView(
            padding: const EdgeInsets.all(16.0),
            children: [
              _buildHeader(invoice),
              const SizedBox(height: 20),
              _buildInfoSection(invoice),
              const SizedBox(height: 20),
              _buildItemsSection(invoice),
              const SizedBox(height: 20),
              _buildTotalSection(invoice),
            ],
          );
        },
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () { /* TODO: Implement PDF Export */ },
        label: const Text('Tải PDF'),
        icon: const Icon(Icons.picture_as_pdf),
      ),
    );
  }

  Widget _buildHeader(Invoice invoice) {
    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            Text(invoice.invoiceNumber, style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 8),
            Chip(
              label: Text(invoice.status.toUpperCase()),
              backgroundColor: invoice.status == 'paid' ? Colors.green[100] : Colors.orange[100],
              labelStyle: TextStyle(color: invoice.status == 'paid' ? Colors.green[800] : Colors.orange[800]),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoSection(Invoice invoice) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Thông tin', style: Theme.of(context).textTheme.titleLarge),
            const Divider(),
            ListTile(
              leading: const Icon(Icons.person),
              title: const Text('Khách hàng'),
              subtitle: Text(invoice.customerName),
            ),
            ListTile(
              leading: const Icon(Icons.calendar_today),
              title: const Text('Ngày phát hành'),
              subtitle: Text(DateFormat('dd/MM/yyyy').format(invoice.issueDate)),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildItemsSection(Invoice invoice) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Chi tiết sản phẩm', style: Theme.of(context).textTheme.titleLarge),
            const Divider(),
            if (invoice.items == null || invoice.items!.isEmpty)
              const Text('Không có sản phẩm.')
            else
              ...invoice.items!.map((item) => ListTile(
                    title: Text('${item.quantity} x ${item.description}'),
                    subtitle: Text('@ ${NumberFormat.decimalPattern().format(item.unitPrice)}'),
                    trailing: Text(NumberFormat.decimalPattern().format(item.totalPrice)),
                  )).toList(),
          ],
        ),
      ),
    );
  }

  Widget _buildTotalSection(Invoice invoice) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            _buildTotalRow('Tổng tiền hàng:', '${NumberFormat.decimalPattern().format(invoice.subtotal)} VND'),
            _buildTotalRow('Thuế (10%):', '${NumberFormat.decimalPattern().format(invoice.tax)} VND'),
            const Divider(),
            _buildTotalRow('TỔNG CỘNG:', '${NumberFormat.decimalPattern().format(invoice.total)} VND', isTotal: true),
          ],
        ),
      ),
    );
  }

  Widget _buildTotalRow(String title, String value, {bool isTotal = false}) {
    final style = TextStyle(
      fontSize: isTotal ? 18 : 16,
      fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [Text(title, style: style), Text(value, style: style)],
      ),
    );
  }
}