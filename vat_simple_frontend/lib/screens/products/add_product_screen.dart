import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vat_simple_frontend/api/product_service.api.dart';
import 'package:vat_simple_frontend/providers/auth_provider.dart';

class AddProductScreen extends StatefulWidget {
  const AddProductScreen({super.key});

  @override
  State<AddProductScreen> createState() => _AddProductScreenState();
}

class _AddProductScreenState extends State<AddProductScreen> {
  final _formKey = GlobalKey<FormState>();
  bool _isLoading = false;
  String _name = '';
  String _description = '';
  double _unitPrice = 0.0;
  String _imageUrl = ''; // <-- Thêm biến cho image url

  Future<void> _saveProduct() async {
    if (!_formKey.currentState!.validate()) return;
    _formKey.currentState!.save();
    setState(() => _isLoading = true);

    final token = Provider.of<AuthProvider>(context, listen: false).token!;
    try {
      // Thêm _imageUrl vào hàm gọi API
      await ProductService().createProduct(token, _name, _description, _unitPrice, _imageUrl);
      if (mounted) Navigator.of(context).pop(true);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(e.toString())));
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Thêm sản phẩm mới')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                decoration: const InputDecoration(labelText: 'Tên sản phẩm *'),
                validator: (v) => (v == null || v.isEmpty) ? 'Không được để trống' : null,
                onSaved: (v) => _name = v!,
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: 'Mô tả'),
                onSaved: (v) => _description = v ?? '',
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: 'Đơn giá *'),
                keyboardType: TextInputType.number,
                validator: (v) {
                  if (v == null || v.isEmpty) return 'Không được để trống';
                  if (double.tryParse(v) == null) return 'Vui lòng nhập số hợp lệ';
                  if (double.parse(v) <= 0) return 'Giá phải lớn hơn 0';
                  return null;
                },
                onSaved: (v) => _unitPrice = double.parse(v!),
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: 'URL Hình ảnh'),
                keyboardType: TextInputType.url,
                onSaved: (v) => _imageUrl = v ?? '',
              ),
              const SizedBox(height: 20),
              if (_isLoading)
                const Center(child: CircularProgressIndicator())
              else
                ElevatedButton(onPressed: _saveProduct, child: const Text('Lưu')),
            ],
          ),
        ),
      ),
    );
  }
}