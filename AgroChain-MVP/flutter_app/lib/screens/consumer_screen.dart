import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_app/api_service.dart';

class ConsumerScreen extends StatefulWidget {
  const ConsumerScreen({super.key});
  @override
  State<ConsumerScreen> createState() => _ConsumerScreenState();
}

class _ConsumerScreenState extends State<ConsumerScreen> {
  final _idController = TextEditingController();
  final ApiService _apiService = ApiService();
  bool _isLoading = false;
  List<dynamic>? _history;
  String? _errorMessage;

  void _traceProduce() async {
    if (_idController.text.isEmpty) return;
    setState(() { _isLoading = true; _history = null; _errorMessage = null; });
    try {
      final historyData = await _apiService.getProduceHistory(_idController.text);
      setState(() => _history = historyData);
    } catch (e) {
      setState(() => _errorMessage = 'Error: $e');
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Trace Produce History')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            TextField(controller: _idController, decoration: const InputDecoration(labelText: 'Enter Produce ID', suffixIcon: Icon(Icons.qr_code))),
            const SizedBox(height: 10),
            SizedBox(width: double.infinity, child: ElevatedButton(onPressed: _traceProduce, child: const Text('Trace'))),
            const SizedBox(height: 20),
            if (_isLoading) const CircularProgressIndicator(),
            if (_errorMessage != null) Text(_errorMessage!, style: const TextStyle(color: Colors.red)),
            if (_history != null)
              Expanded(
                child: ListView.builder(
                  itemCount: _history!.length,
                  itemBuilder: (context, index) {
                    final item = jsonDecode(_history![index]);
                    return Card(
                      margin: const EdgeInsets.symmetric(vertical: 8.0),
                      child: ListTile(
                        leading: CircleAvatar(child: Text('${index + 1}')),
                        title: Text('Owner: ${item['owner']}'),
                        subtitle: Text('Crop: ${item['crop']} | Qty: ${item['quantity']} kg'),
                        trailing: Text(DateTime.parse(item['timestamp']).toLocal().toString().substring(0, 16)),
                      ),
                    );
                  },
                ),
              ),
          ],
        ),
      ),
    );
  }
}