import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/screens/customers/customer_list_screen.dart';
import 'package:vat_simple_frontend/screens/products/product_list_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    // Dùng Consumer để lắng nghe và lấy dữ liệu user từ AuthProvider
    return Consumer<AuthProvider>(
      builder: (context, authProvider, child) {
        // Nếu user chưa được tải, hiển thị vòng xoay
        if (authProvider.user == null) {
          return const Scaffold(
            body: Center(
              child: CircularProgressIndicator(),
            ),
          );
        }

        // Nếu user đã được tải, hiển thị giao diện chính
        return Scaffold(
          appBar: AppBar(
            // Hiển thị tên công ty trên AppBar
            title: Text(authProvider.user!.companyName),
            centerTitle: false,
            actions: [
              // Nút đăng xuất
              IconButton(
                icon: const Icon(Icons.logout),
                tooltip: 'Đăng xuất',
                onPressed: () {
                  Provider.of<AuthProvider>(context, listen: false).logout();
                },
              ),
            ],
          ),
          body: ListView(
            padding: const EdgeInsets.all(16.0),
            children: [
              // --- Phần chào mừng ---
              Card(
                elevation: 2,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Xin chào,',
                        style: TextStyle(fontSize: 16, color: Colors.grey),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        authProvider.user!.email,
                        style: const TextStyle(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              // --- Phần chức năng ---
              Text(
                'Bảng điều khiển',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              Card(
                elevation: 2,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                clipBehavior: Clip.antiAlias, // Để bo góc cho ListTile
                child: Column(
                  children: [
                    _buildFeatureTile(
                      context: context,
                      icon: Icons.people_outline,
                      title: 'Quản lý Khách hàng',
                      onTap: () {
                        Navigator.of(context).push(MaterialPageRoute(
                          builder: (context) => const CustomerListScreen(),
                        ));
                      },
                    ),
                    const Divider(height: 1, indent: 16, endIndent: 16),
                    _buildFeatureTile(
                      context: context,
                      icon: Icons.inventory_2_outlined,
                      title: 'Quản lý Sản phẩm',
                      onTap: () {
                        Navigator.of(context).push(MaterialPageRoute(
                          builder: (context) => const ProductListScreen(),
                        ));
                      },
                    ),
                    const Divider(height: 1, indent: 16, endIndent: 16),
                    _buildFeatureTile(
                      context: context,
                      icon: Icons.receipt_long_outlined,
                      title: 'Tạo Hóa đơn mới',
                      onTap: () {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(content: Text('Chức năng sẽ được phát triển!'))
                        );
                      },
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  // Widget helper để tạo các ô chức năng cho đồng bộ
  Widget _buildFeatureTile({
    required BuildContext context,
    required IconData icon,
    required String title,
    required VoidCallback onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: Theme.of(context).primaryColor),
      title: Text(title, style: const TextStyle(fontWeight: FontWeight.w500)),
      trailing: const Icon(Icons.arrow_forward_ios, size: 16),
      onTap: onTap,
    );
  }
}