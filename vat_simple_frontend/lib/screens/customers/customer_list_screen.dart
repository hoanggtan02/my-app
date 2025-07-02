import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/customer_service.api.dart';
import 'package:vat_simple_frontend/models/customer.model.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';

class CustomerListScreen extends StatefulWidget {
  const CustomerListScreen({super.key});

  @override
  State<CustomerListScreen> createState() => _CustomerListScreenState();
}

class _CustomerListScreenState extends State<CustomerListScreen> {
  late Future<List<Customer>> _customersFuture;
  final CustomerService _customerService = CustomerService();

  @override
  void initState() {
    super.initState();
    // Lấy token từ AuthProvider và gọi API
    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    _customersFuture = _customerService.getCustomers(token);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Danh sách Khách hàng'),
      ),
      body: FutureBuilder<List<Customer>>(
        future: _customersFuture,
        builder: (context, snapshot) {
          // Trong khi đang tải dữ liệu
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          // Nếu có lỗi
          if (snapshot.hasError) {
            return Center(child: Text('Lỗi: ${snapshot.error}'));
          }
          // Nếu không có dữ liệu
          if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('Chưa có khách hàng nào.'));
          }

          // Hiển thị danh sách
          final customers = snapshot.data!;
          return ListView.builder(
            itemCount: customers.length,
            itemBuilder: (ctx, index) {
              final customer = customers[index];
              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                child: ListTile(
                  leading: CircleAvatar(child: Text(customer.name[0])),
                  title: Text(customer.name),
                  subtitle: Text(customer.taxCode ?? 'Chưa có MST'),
                  trailing: const Icon(Icons.arrow_forward_ios),
                  onTap: () {
                    // Xử lý khi nhấn vào một khách hàng
                  },
                ),
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Xử lý khi nhấn nút thêm khách hàng mới
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}