import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/screens/home/home_screen.dart';
import 'package:vat_simple_frontend/screens/login/login_screen.dart';

class AuthWrapper extends StatefulWidget {
  const AuthWrapper({super.key});

  @override
  State<AuthWrapper> createState() => _AuthWrapperState();
}

class _AuthWrapperState extends State<AuthWrapper> {
  // Biến để theo dõi trạng thái loading ban đầu
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    // Gọi hàm tryAutoLogin MỘT LẦN DUY NHẤT ở đây
    Provider.of<AuthProvider>(context, listen: false).tryAutoLogin().then((_) {
      // Sau khi kiểm tra xong, cập nhật lại trạng thái để tắt loading
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    // Nếu đang loading, hiển thị vòng xoay
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    // Lắng nghe sự thay đổi để build lại UI khi login/logout
    final authProvider = Provider.of<AuthProvider>(context);

    // Dựa vào trạng thái isAuth để quyết định hiển thị màn hình nào
    if (authProvider.isAuth) {
      return const HomeScreen();
    } else {
      return const LoginScreen();
    }
  }
}