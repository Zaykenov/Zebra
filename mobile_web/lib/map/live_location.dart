import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_smart_dialog/flutter_smart_dialog.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';
import 'package:maps_launcher/maps_launcher.dart';
import 'package:mobile_web/map/map_style.dart';
import 'package:permission_handler/permission_handler.dart';

import '../helpers/api_caller.dart';

class MapScreen extends StatefulWidget {
  @override
  _MapScreenState createState() => _MapScreenState();
}

class _MapScreenState extends State<MapScreen> {
  GoogleMapController? _mapController;
  LatLng? _currentPosition;
  Set<Marker> zebraLocationMarkers = {}; // Change to nullable type

  @override
  void initState() {
    super.initState();
    _getCurrentLocation();
    startLoadZebraLocations();
  }

  void _getCurrentLocation() async {
    // Check if location permission is granted
    if (await Permission.location.isGranted) {
      Position position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.best,
        timeLimit: const Duration(seconds: 10),
      );

      setState(() {
        _currentPosition = LatLng(position.latitude, position.longitude);
      });
    } else {
      // Request location permission
      var status = await Permission.location.request();
      if (status.isGranted) {
        // If permission is granted, get the current location
        _getCurrentLocation();
      } else if (status.isDenied) {
        // Handle the case when the user denied the permission
        SmartDialog.show(builder: (context) {
          return Container(
            height: 150,
            width: 300,
            decoration: BoxDecoration(
              color: const Color(0xfffafafa),
              borderRadius: BorderRadius.circular(10),
            ),
            alignment: Alignment.center,
            child:
                Column(mainAxisAlignment: MainAxisAlignment.center, children: [
              const Padding(
                padding: EdgeInsets.only(top: 10),
                child: Text('Location Permission is Denied',
                    style:
                        TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
              ),
              const SizedBox(height: 20),
              TextButton(
                  onPressed: () {
                    SmartDialog.dismiss();
                  },
                  child: const Text("Close")),
            ]),
          );
        });
      } else if (status.isPermanentlyDenied) {
        // Handle the case when the user permanently denied the permission
        SmartDialog.show(builder: (context) {
          return Container(
            height: 150,
            width: 300,
            decoration: BoxDecoration(
              color: const Color(0xfffafafa),
              borderRadius: BorderRadius.circular(10),
            ),
            alignment: Alignment.center,
            child:
                Column(mainAxisAlignment: MainAxisAlignment.center, children: [
              const Padding(
                padding: EdgeInsets.only(top: 10),
                child: Text('Location Permission is Permanently Denied',
                    style:
                        TextStyle(fontWeight: FontWeight.bold, fontSize: 18)),
              ),
              const SizedBox(height: 20),
              TextButton(
                  onPressed: () {
                    SmartDialog.dismiss();
                  },
                  child: const Text("Close")),
            ]),
          );
        });
      }
    }
  }

  void _centerOnCurrentLocation() async {
    if (_mapController != null && _currentPosition != null) {
      // Check for null
      _mapController!.animateCamera(
        CameraUpdate.newCameraPosition(
          CameraPosition(
            target: _currentPosition!,
            zoom: 15.0,
          ),
        ),
      );
    }
  }

  Future<void> startLoadZebraLocations() async {
    try {
      final response = await ApiCaller().fetchZebraLocations();

      BitmapDescriptor customIcon = await BitmapDescriptor.fromAssetImage(
        const ImageConfiguration(),
        'assets/map_logo.png',
      );

      Set<Marker> markers = response
          .map<Marker>(
            (e) => Marker(
              markerId: MarkerId(e['ShopId'].toString()),
              position: LatLng(e['Latitude'], e['Longitude']),
              icon: customIcon,
              onTap: () {
                SmartDialog.show(builder: (context) {
                  return Container(
                    height: 150,
                    width: 300,
                    decoration: BoxDecoration(
                      color: const Color(0xfffafafa),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    alignment: Alignment.center,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(top: 10),
                          child: Text(e['Address'],
                              style: const TextStyle(
                                  fontWeight: FontWeight.bold, fontSize: 18)),
                        ),
                        const SizedBox(height: 20),
                        TextButton(
                          style: ButtonStyle(
                            foregroundColor:
                                MaterialStateProperty.all<Color>(Colors.blue),
                          ),
                          onPressed: () async {
                            _launchMaps(
                              e['Latitude'],
                              e['Longitude'],
                              e['Name'],
                            );
                          },
                          child: const Text(
                            'Проложить маршрут',
                            style: TextStyle(color: Colors.black),
                          ),
                        )
                      ],
                    ),
                  );
                });
              },
            ),
          )
          .toSet();

      setState(() {
        zebraLocationMarkers = markers;
      });
    } catch (error) {
      if (kDebugMode) {
        print('Error loading shop locations: $error');
      }
      // Handle the error as needed
    }
  }

  void _launchMaps(double latitude, double longitude, String name) {
    MapsLauncher.launchCoordinates(latitude, longitude, name);
  }

  Marker? _findNearestMarker() {
    if (_currentPosition == null || zebraLocationMarkers.isEmpty) return null;

    double minDistance = double.infinity;
    Marker? nearestMarker; // Initialize nearestMarker as nullable

    for (Marker marker in zebraLocationMarkers) {
      double distance = Geolocator.distanceBetween(
        _currentPosition!.latitude,
        _currentPosition!.longitude,
        marker.position.latitude,
        marker.position.longitude,
      );

      if (distance < minDistance) {
        minDistance = distance;
        nearestMarker = marker;
      }
    }

    return nearestMarker;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Карта заведений')),
      body: _currentPosition == null
          ? const Center(
              child: CircularProgressIndicator(),
            )
          : Stack(
              children: [
                GoogleMap(
                  markers: zebraLocationMarkers,
                  rotateGesturesEnabled: true,
                  initialCameraPosition: CameraPosition(
                    target: _currentPosition!,
                    zoom: 15.0,
                  ),
                  onMapCreated: (GoogleMapController controller) {
                    _mapController = controller;
                    _mapController?.setMapStyle(MapStyle().aubergine);
                  },
                  myLocationEnabled: true,
                  myLocationButtonEnabled: false, // Hide the location button
                  compassEnabled: true,
                  zoomControlsEnabled: false, // Show zoom in/out buttons
                ),
                Positioned(
                  bottom: 16,
                  left: 0,
                  right: 0,
                  child: Center(
                    child: SizedBox(
                      width: 200, // Adjust the width as per your requirement
                      height: 55, // Adjust the height as per your requirement
                      child: ElevatedButton(
                        onPressed: () {
                          Marker? nearestMarker = _findNearestMarker();
                          if (nearestMarker != null) {
                            _mapController?.animateCamera(
                              CameraUpdate.newCameraPosition(
                                CameraPosition(
                                  target: nearestMarker.position,
                                  zoom: 16.0,
                                ),
                              ),
                            );
                          } else {
                            // Handle the case when no markers are available
                            if (kDebugMode) {
                              print('No markers found.');
                            }
                          }
                        },
                        child: const Text('Найти ZebraCoffee',
                            style: TextStyle(
                                fontSize: 15)), // Adjust font size as needed
                      ),
                    ),
                  ),
                ),
              ],
            ),
      floatingActionButton: _currentPosition == null
          ? null // Disable the button until _currentPosition is initialized
          : FloatingActionButton(
              onPressed: _centerOnCurrentLocation,
              tooltip: 'Center on Current Location',
              child: const Icon(Icons.location_searching),
            ),
    );
  }
}
