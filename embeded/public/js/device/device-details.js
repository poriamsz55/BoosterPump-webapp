// get all parts that belong to this device
// render them in the device details

const urlParams = new URLSearchParams(window.location.search);
const deviceId = urlParams.get('id');

// get deviceParts from local storage
const deviceParts = JSON.parse(localStorage.getItem('deviceParts')) || [];

// get deviceParts from database
const getDeviceParts = async () => {
    const response = await fetch(`${HTTP_URL}/devicePart/getByDeviceId?id=${deviceId}`);
    const deviceParts = await response.json();
    return deviceParts;
}

// render deviceParts in device details
const renderDeviceParts = async () => {
    const deviceParts = await getDeviceParts();
    const devicePartList = document.getElementById('devicePartList');
    devicePartList.innerHTML = '';
    deviceParts.forEach(devicePart => {
        const li = document.createElement('li');
        li.textContent = `${devicePart.partName} x${devicePart.count}`;
        devicePartList.appendChild(li);
    });
}

// addPartToDeviceBtn
const addPartToDeviceBtn = document.getElementById('addPartToDeviceBtn');
// get all parts from database
const getPartsFromDB = async () => {
    const response = await fetch(`${HTTP_URL}/part/getAll`);
    const parts = await response.json();
    return parts;
}

// render parts in modal
const renderParts = async (parts) => {
    const partsGrid = document.getElementById('partsGrid');
    partsGrid.innerHTML = '';
    parts.forEach(part => {
        // check if part already exists in deviceParts
        const partExists = deviceParts.some(devicePart => devicePart.partId === part.id);
        if (partExists) {
           // TODO
        }
        const partCard = document.createElement('div');
        partCard.classList.add('part-card');
        partCard.innerHTML = `
            <div class="part-name">${part.name}</div>
            <div class="part-count">${part.count}</div>
        `;
        partsGrid.appendChild(partCard);
    });
}
// open modal and add disabled class to partsGrid that already exists in deviceParts
const openModalAndRenderParts = async () => {
    openModal();
    const parts = await getPartsFromDB();
    renderParts(parts);
}

// open modal
const openModal = () => {
    document.getElementById('addPartToDeviceModal').style.display = 'block';
}

// close modal
const closeModal = () => {
    document.getElementById('addPartToDeviceModal').style.display = 'none';
}

renderDeviceParts();