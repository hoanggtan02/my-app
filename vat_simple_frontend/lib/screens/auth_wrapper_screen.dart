import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/screens/login/login_screen.dart';
import 'package:vat_simple_frontend/screens/main_screen.dart'; 

class AuthWrapper extends StatefulWidget {
  const AuthWrapper({super.key});

  @override
  State<AuthWrapper> createState() => _AuthWrapperState();
}

class _AuthWrapperState extends State<AuthWrapper> {
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    Provider.of<AuthProvider>(context, listen: false).tryAutoLogin().then((_) {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    final authProvider = Provider.of<AuthProvider>(context);

    if (authProvider.isAuth) {
      return const MainScreen(); 
    } 
    else {
      return const LoginScreen();
    }
  }
}