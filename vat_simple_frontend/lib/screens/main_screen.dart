import 'package:flutter/material.dart';
import 'package:vat_simple_frontend/screens/customers/customer_list_screen.dart';
import 'package:vat_simple_frontend/screens/home/home_screen.dart';
import 'package:vat_simple_frontend/screens/products/product_list_screen.dart';
import 'package:vat_simple_frontend/screens/invoices/invoice_list_screen.dart'; 
import 'package:vat_simple_frontend/screens/reports/report_screen.dart'; 



class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  int _selectedIndex = 0;

  // Danh sách các màn hình chính của bạn
  static const List<Widget> _widgetOptions = <Widget>[
    HomeScreen(),
    CustomerListScreen(),
    ProductListScreen(),
    InvoiceListScreen(), 
    ReportScreen(), 
    Center(child: Text('Lịch sử Hóa đơn')), 
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    // Dùng LayoutBuilder để kiểm tra chiều rộng màn hình
    return LayoutBuilder(
      builder: (context, constraints) {
        // Nếu màn hình rộng (desktop/web) -> Dùng menu dọc
        if (constraints.maxWidth > 600) {
          return Scaffold(
            body: Row(
              children: <Widget>[
                NavigationRail(
                  selectedIndex: _selectedIndex,
                  onDestinationSelected: _onItemTapped,
                  labelType: NavigationRailLabelType.all,
                  extended: constraints.maxWidth > 800, // Tự động mở rộng menu trên màn hình rất rộng
                  destinations: const <NavigationRailDestination>[
                    NavigationRailDestination(
                      icon: Icon(Icons.home_outlined),
                      selectedIcon: Icon(Icons.home),
                      label: Text('Trang chủ'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.people_outline),
                      selectedIcon: Icon(Icons.people),
                      label: Text('Khách hàng'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.inventory_2_outlined),
                      selectedIcon: Icon(Icons.inventory_2),
                      label: Text('Sản phẩm'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.receipt_long_outlined),
                      selectedIcon: Icon(Icons.receipt_long),
                      label: Text('Hóa đơn'),
                    ),
                    NavigationRailDestination(
                      icon: Icon(Icons.bar_chart_outlined),
                      selectedIcon: Icon(Icons.bar_chart),
                      label: Text('Báo cáo'),
                    ),
                  ],
                ),
                const VerticalDivider(thickness: 1, width: 1),
                // Hiển thị nội dung trang được chọn
                Expanded(
                  child: _widgetOptions.elementAt(_selectedIndex),
                ),
              ],
            ),
          );
        }
        // Nếu màn hình hẹp (điện thoại) -> Dùng thanh điều hướng dưới đáy
        else {
          return Scaffold(
            body: Center(
              child: _widgetOptions.elementAt(_selectedIndex),
            ),
            bottomNavigationBar: BottomNavigationBar(
              items: const <BottomNavigationBarItem>[
                BottomNavigationBarItem(
                  icon: Icon(Icons.home),
                  label: 'Trang chủ',
                ),
                BottomNavigationBarItem(
                  icon: Icon(Icons.people),
                  label: 'Khách hàng',
                ),
                BottomNavigationBarItem(
                  icon: Icon(Icons.inventory_2),
                  label: 'Sản phẩm',
                ),
                BottomNavigationBarItem(
                  icon: Icon(Icons.receipt_long),
                  label: 'Hóa đơn',
                ),
                BottomNavigationBarItem(
                  icon: Icon(Icons.bar_chart),
                  label: 'Báo cáo',
                ),
              ],
              currentIndex: _selectedIndex,
              selectedItemColor: Theme.of(context).primaryColor,
              unselectedItemColor: Colors.grey,
              showUnselectedLabels: true,
              onTap: _onItemTapped,
            ),
          );
        }
      },
    );
  }
}