import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';
import 'package:vat_simple_frontend/screens/auth_wrapper_screen.dart';
import 'package:vat_simple_frontend/providers/cart_provider.dart';

void main() {
  runApp(const MyApp());
}

// Trong file main.dart
class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Dùng MultiProvider để cung cấp nhiều provider
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (ctx) => AuthProvider()),
        ChangeNotifierProvider(create: (ctx) => CartProvider()), 
      ],
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'VAT Simple App',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigo),
          useMaterial3: true,
          scaffoldBackgroundColor: Colors.grey[100],
          appBarTheme: const AppBarTheme(
            backgroundColor: Colors.white,
            foregroundColor: Colors.black87,
            elevation: 1,
          ),
          cardTheme: CardTheme(
            elevation: 1,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          ),
        ),
        home: const AuthWrapper(),
      ),
    );
  }
}