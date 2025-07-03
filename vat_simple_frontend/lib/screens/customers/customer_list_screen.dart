import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/customer_service.api.dart';
import 'package:vat_simple_frontend/models/customer.model.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/screens/customers/add_edit_customer_screen.dart';

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
    _refreshCustomers(); 
  }

  Future<void> _refreshCustomers() async {
    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    setState(() {
      _customersFuture = _customerService.getCustomers(token);
    });
  }

  // Hàm điều hướng và chờ kết quả để làm mới danh sách
  void _navigateAndRefresh(Widget screen) async {
    final result = await Navigator.of(context).push<bool>(
      MaterialPageRoute(builder: (context) => screen),
    );
    // Nếu màn hình con trả về true (tức là đã có thay đổi), làm mới danh sách
    if (result == true && mounted) {
      _refreshCustomers();
    }
  }

  // Hàm xử lý xóa khách hàng
  void _deleteCustomer(String customerId) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Xác nhận Xóa'),
        content: const Text('Bạn có chắc chắn muốn xóa khách hàng này không?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('Hủy'),
          ),
          TextButton(
            onPressed: () async {
              Navigator.of(ctx).pop();
              try {
                final token = Provider.of<AuthProvider>(context, listen: false).token!;
                await _customerService.deleteCustomer(token, customerId);
                _refreshCustomers(); // Làm mới danh sách sau khi xóa
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Lỗi khi xóa: ${e.toString()}')),
                  );
                }
              }
            },
            child: const Text('Xóa', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Danh sách Khách hàng'),
      ),
      body: RefreshIndicator(
        onRefresh: _refreshCustomers,
        child: FutureBuilder<List<Customer>>(
          future: _customersFuture,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            }
            if (snapshot.hasError) {
              return Center(child: Text('Lỗi: ${snapshot.error}'));
            }
            if (!snapshot.hasData || snapshot.data!.isEmpty) {
              return const Center(
                child: Text(
                  'Chưa có khách hàng nào.',
                  style: TextStyle(fontSize: 18, color: Colors.grey),
                ),
              );
            }

            final customers = snapshot.data!;
            return ListView.builder(
              itemCount: customers.length,
              itemBuilder: (ctx, index) {
                final customer = customers[index];
                return Card(
                  margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  child: ListTile(
                    leading: CircleAvatar(child: Text(customer.name[0].toUpperCase())),
                    title: Text(customer.name),
                    subtitle: Text(customer.taxCode ?? 'Chưa có MST'),
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        IconButton(
                          icon: Icon(Icons.edit, color: Colors.blue[400]),
                          tooltip: 'Sửa',
                          onPressed: () => _navigateAndRefresh(AddEditCustomerScreen(customer: customer)),
                        ),
                        IconButton(
                          icon: Icon(Icons.delete, color: Colors.red[400]),
                          tooltip: 'Xóa',
                          onPressed: () => _deleteCustomer(customer.id),
                        ),
                      ],
                    ),
                  ),
                );
              },
            );
          },
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _navigateAndRefresh(const AddEditCustomerScreen()),
        tooltip: 'Thêm khách hàng mới',
        child: const Icon(Icons.add),
      ),
    );
  }
}