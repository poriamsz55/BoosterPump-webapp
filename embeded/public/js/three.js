 // Three.js Scene Setup
 const scene = new THREE.Scene();
 const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
 const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
 
 renderer.setSize(window.innerWidth, window.innerHeight);
 document.getElementById('scene-container').appendChild(renderer.domElement);

 // Particles
 const particlesGeometry = new THREE.BufferGeometry();
 const particlesCount = 1500;
 const posArray = new Float32Array(particlesCount * 3);

 for(let i = 0; i < particlesCount * 3; i++) {
     posArray[i] = (Math.random() - 0.5) * 5;
 }

 particlesGeometry.setAttribute('position', new THREE.BufferAttribute(posArray, 3));

 const particlesMaterial = new THREE.PointsMaterial({
     size: 0.005,
     color: '#4834d4',
     transparent: true,
     opacity: 0.8,
 });

 const particlesMesh = new THREE.Points(particlesGeometry, particlesMaterial);
//  scene.add(particlesMesh);

 camera.position.z = 2;

 // Animation
 function animate() {
     requestAnimationFrame(animate);
     particlesMesh.rotation.y += 0.001;
     particlesMesh.rotation.x += 0.001;
     renderer.render(scene, camera);
 }

 // Handle Window Resize
 window.addEventListener('resize', () => {
     camera.aspect = window.innerWidth / window.innerHeight;
     camera.updateProjectionMatrix();
     renderer.setSize(window.innerWidth, window.innerHeight);
 });

 // Start Animation
 animate();
