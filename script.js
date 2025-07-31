// API_URL - адрес сервиса API (замените на свой для production, если надо)
const API_URL = "http://localhost:8080/api"; // или через nginx proxy

async function loadApod() {
    const apodResult = document.getElementById('apod-result');
    apodResult.innerHTML = '<p class="loading">Загрузка фото дня...</p>';
    try {
        const response = await fetch(`${API_URL}/apod`);
        if (!response.ok) throw new Error("Ошибка сети");
        const data = await response.json();
        apodResult.innerHTML = `
            <h3>${data.title}</h3>
            <p>${data.date}</p>
            <img src="${data.url}" alt="APOD">
            <p>${data.explanation || ''}</p>
        `;
    } catch (error) {
        apodResult.innerHTML = '<p>Ошибка при загрузке фото дня.</p>';
    }
}

async function loadRoverPhotos(sol, camera) {
    const roverResult = document.getElementById('rover-result');
    roverResult.innerHTML = '<p class="loading">Загрузка фотографий марсохода...</p>';
    try {
        const response = await fetch(`${API_URL}/rover?sol=${sol}&camera=${camera}`);
        if (!response.ok) throw new Error("Ошибка сети");
        const data = await response.json();
        roverResult.innerHTML = '';
        if (!data.photos || data.photos.length === 0) {
            roverResult.innerHTML = '<p>Фотографии не найдены для выбранного Sol и камеры.</p>';
            return;
        }
        data.photos.forEach(photo => {
            const photoDiv = document.createElement('div');
            photoDiv.classList.add('photo');
            photoDiv.innerHTML = `<img src="${photo.img_src}" alt="Mars Rover Photo"><p>Дата: ${photo.earth_date}</p>`;
            roverResult.appendChild(photoDiv);
        });
    } catch (error) {
        roverResult.innerHTML = '<p>Ошибка при загрузке фотографий.</p>';
    }
}

document.getElementById('rover-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const sol = document.getElementById('sol').value;
    const camera = document.getElementById('camera').value;
    loadRoverPhotos(sol, camera);
});

function toggleTheme() {
    document.body.classList.toggle('dark');
}
document.getElementById('theme-toggle').addEventListener('click', toggleTheme);

// Инициализация
loadApod();