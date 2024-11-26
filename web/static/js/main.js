let modelConfigs = {};
document.getElementById('optimizer').addEventListener('change', function(e) {
    const container = document.getElementById('trainable_params_container');
    container.style.display = e.target.value ? 'block' : 'none';
});
function updateFormFields(modelName) {
    const form = document.getElementById('calculatorForm');
    if (!modelName || !modelConfigs[modelName]) {
        form.reset();
        return;
    }
    const config = modelConfigs[modelName];
    console.log("Selected model config:", config);
    
    document.getElementById('model_size').value = config.model_size || '';
    document.getElementById('hidden_size').value = config.hidden_size || '';
    document.getElementById('num_hidden_layers').value = config.num_hidden_layers || '';
    document.getElementById('num_attention_heads').value = config.num_attention_heads || '';
    document.getElementById('sequence_length').value = config.max_position_embeddings || 4096;
    document.getElementById('batch_size').value = 1;
    
    const dtypeSelect = document.getElementById('torch_dtype');
    const dtype = config.torch_dtype || 'float16';
    for (let i = 0; i < dtypeSelect.options.length; i++) {
        if (dtypeSelect.options[i].value === dtype) {
            dtypeSelect.selectedIndex = i;
            break;
        }
    }
}

function displayGPURecommendations(recommendations) {
    console.log("Displaying GPU recommendations:", recommendations);
    if (!recommendations || !Array.isArray(recommendations) || recommendations.length === 0) {
        console.warn("No GPU recommendations available");
        return;
    }
    const gpuCardsDiv = document.querySelector('.gpu-cards');
    if (!gpuCardsDiv) {
        console.error("GPU cards container not found");
        return;
    }
    gpuCardsDiv.innerHTML = '';
    const sortedRecs = recommendations.sort((a, b) => {
        if (a.training_supported !== b.training_supported) {
            return b.training_supported - a.training_supported;
        }
        return a.utilization_score - b.utilization_score;
    });
    const gpuSections = {
        training: {
            title: 'Recommended for Training',
            gpus: sortedRecs.filter(rec => rec.training_supported && rec.utilization_score < 95)
        },
        inference: {
            title: 'Recommended for Inference',
            gpus: sortedRecs.filter(rec => !rec.training_supported && rec.utilization_score < 95)
        }
    };
    Object.values(gpuSections).forEach(section => {
        if (section.gpus.length === 0) return;

        const sectionDiv = document.createElement('div');
        sectionDiv.className = 'gpu-section';
        sectionDiv.innerHTML = `
            <h4 class="gpu-section-title">${section.title}</h4>
            <div class="gpu-section-cards">
                ${section.gpus.map(rec => {
                    const utilizationClass = rec.utilization_score > 90 ? 'high-utilization' :
                                          rec.utilization_score > 70 ? 'medium-utilization' : 'optimal-utilization';
                    
                    return `
                        <div class="gpu-card ${utilizationClass}">
                            <div class="gpu-usage-type">Memory Utilization: ${rec.utilization_score.toFixed(1)}%</div>
                            <h4>${rec.gpu.name || 'Unknown GPU'}</h4>
                            <div class="gpu-specs">
                                <div class="gpu-spec">
                                    <span class="gpu-spec-label">Memory:</span>
                                    <span class="gpu-spec-value">${rec.gpu.memory_gb || 0} GB</span>
                                </div>
                                <div class="gpu-spec">
                                    <span class="gpu-spec-label">Bandwidth:</span>
                                    <span class="gpu-spec-value">${rec.gpu.bandwidth_tbs || 0} TB/s</span>
                                </div>
                                <div class="gpu-spec">
                                    <span class="gpu-spec-label">Price:</span>
                                    <span class="gpu-spec-value">$${rec.gpu.price_usd || 0}</span>
                                </div>
                                <div class="gpu-spec">
                                    <span class="gpu-spec-label">Required:</span>
                                    <span class="gpu-spec-value">${rec.num_gpus || 1} GPU(s)</span>
                                </div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
        `;
        gpuCardsDiv.appendChild(sectionDiv);
    });

    const gpuContainer = document.querySelector('.gpu-recommendations-container');
    if (gpuContainer) {
        gpuContainer.style.display = 'block';
    }
}

function displayResults(data) {
    const resultsContainer = document.querySelector('.results-container');
    const gpuContainer = document.querySelector('.gpu-recommendations-container');
    resultsContainer.style.display = 'block';

    resultsContainer.innerHTML = `
        <div class="memory-sections">
            <!-- Inference Memory Section -->
            <div class="memory-section">
                <h3>Inference Memory Requirements</h3>
                <div class="memory-breakdown">
                    <div class="memory-group">
                        <div class="memory-item">
                            <span class="memory-label">Model Weights:</span>
                            <span class="memory-value">${data.model_weights}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">KV Cache:</span>
                            <span class="memory-value">${data.kv_cache}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">Activation Memory:</span>
                            <span class="memory-value">${data.activation_memory}</span>
                        </div>
                    </div>
                    <div class="memory-total">
                        <span class="memory-label">Total Inference Memory:</span>
                        <span class="memory-value">${data.inference_memory}</span>
                    </div>
                </div>
            </div>

            ${data.training_memory ? `
            <!-- Training Memory Section -->
            <div class="memory-section">
                <h3>Training Memory Requirements</h3>
                <div class="memory-breakdown">
                    <div class="memory-group">
                        <div class="memory-item">
                            <span class="memory-label">Model Weights:</span>
                            <span class="memory-value">${data.model_weights}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">KV Cache:</span>
                            <span class="memory-value">${data.kv_cache}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">Activation Memory:</span>
                            <span class="memory-value">${data.activation_memory}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">Optimizer Memory:</span>
                            <span class="memory-value">${data.optimizer_memory}</span>
                        </div>
                        <div class="memory-item">
                            <span class="memory-label">Gradients Memory:</span>
                            <span class="memory-value">${data.gradients_memory}</span>
                        </div>
                    </div>
                    <div class="memory-total">
                        <span class="memory-label">Total Training Memory:</span>
                        <span class="memory-value">${data.training_memory}</span>
                    </div>
                </div>
            </div>
            ` : ''}
        </div>
    `;


    gpuContainer.innerHTML = '<h2>GPU Recommendations</h2>';
    if (data.inference_gpus && data.inference_gpus.length > 0) {
        const inferenceSection = document.createElement('div');
        inferenceSection.innerHTML = `
            <div class="gpu-section-cards">
                ${data.inference_gpus.map((gpu, index) => createGPUCard(gpu, 'Inference', index === 0)).join('')}
            </div>
        `;
        gpuContainer.appendChild(inferenceSection);
    }

    if (data.training_gpus && data.training_gpus.length > 0) {
        const trainingSection = document.createElement('div');
        trainingSection.innerHTML = `
            <div class="gpu-section-cards">
                ${data.training_gpus.map((gpu, index) => createGPUCard(gpu, 'Training', index === 0)).join('')}
            </div>
        `;
        gpuContainer.appendChild(trainingSection);
    }
}

function createGPUCard(rec, type, isFirst) {
    const utilizationClass = getUtilizationClass(rec.utilization_score / 100);
    const totalCost = (rec.gpu.price_usd * rec.num_gpus).toFixed(2);
    
    return `
        <div class="gpu-card ${utilizationClass} ${isFirst ? 'recommended' : ''}">
            <span class="gpu-usage-type">${type}</span>
            <h4>${rec.gpu.name} ${rec.num_gpus > 1 ? `(${rec.num_gpus}x)` : ''}</h4>
            <div class="gpu-specs">
                <div class="gpu-spec">
                    <span class="gpu-spec-label">Memory Per GPU</span>
                    <span class="gpu-spec-value">${rec.gpu.memory_gb} GB</span>
                </div>
                <div class="gpu-spec">
                    <span class="gpu-spec-label">Memory Utilization</span>
                    <span class="gpu-spec-value">${(rec.utilization_score).toFixed(1)}%</span>
                </div>
                <div class="gpu-spec">
                    <span class="gpu-spec-label">Bandwidth</span>
                    <span class="gpu-spec-value">${rec.gpu.bandwidth_tbs.toFixed(2)} TB/s</span>
                </div>
                <div class="gpu-spec">
                    <span class="gpu-spec-label">Cost Per GPU</span>
                    <span class="gpu-spec-value">$${rec.gpu.price_usd.toFixed(2)}</span>
                </div>
                ${rec.num_gpus > 1 ? `
                <div class="gpu-spec total-cost">
                    <span class="gpu-spec-label">Total Cost (${rec.num_gpus}x)</span>
                    <span class="gpu-spec-value">$${totalCost}</span>
                </div>
                ` : ''}
            </div>
        </div>
    `;
}

function getUtilizationClass(utilization) {
    if (utilization > 0.9) return 'high-utilization';
    if (utilization > 0.7) return 'medium-utilization';
    return 'optimal-utilization';
}


document.querySelector('form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    try {
        const resultsContainer = document.querySelector('.results-container');
        resultsContainer.innerHTML = '<div class="loading">Calculating...</div>';
        const formData = new FormData(this);
        const data = {};
        
        // Convert form data to proper types
        data.model_size = parseFloat(formData.get('model_size') || '0');
        data.hidden_size = parseInt(formData.get('hidden_size') || '0', 10);
        data.num_hidden_layers = parseInt(formData.get('num_hidden_layers') || '0', 10);
        data.num_attention_heads = parseInt(formData.get('num_attention_heads') || '0', 10);
        data.sequence_length = parseInt(formData.get('sequence_length') || '0', 10);
        data.batch_size = parseInt(formData.get('batch_size') || '0', 10);
        data.torch_dtype = formData.get('torch_dtype') || 'float32';
        data.optimizer = formData.get('optimizer') || '';

        console.log("Sending calculation request:", data);

        // Make API request
        const response = await fetch('/api/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || `HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        console.log("API Response:", result);
        const isTraining = data.optimizer && data.optimizer !== '';
        displayResults(result);
        
    } catch (error) {
        console.error('Error:', error);
        const resultsContainer = document.querySelector('.results-container');
        resultsContainer.innerHTML = `
            <div class="error">
                <h3>Error</h3>
                <p>${error.message || 'An error occurred while calculating memory requirements.'}</p>
                <p>Please check your inputs and try again.</p>
            </div>
        `;
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const modelSelect = document.getElementById('model_select');
    const modelsScript = document.querySelector('script[data-models]');
    if (modelsScript) {
        modelConfigs = JSON.parse(modelsScript.getAttribute('data-models'));
        console.log("Loaded model configs:", modelConfigs);
    }
    modelSelect.addEventListener('change', (e) => {
        updateFormFields(e.target.value);
    });
});

function showDocumentation() {
    fetch('/documentation')
        .then(response => response.text())
        .then(markdown => {
            const modal = document.getElementById('documentationModal');
            const content = document.getElementById('documentationContent');
            content.innerHTML = marked.parse(markdown);
            modal.style.display = 'block';
        })
        .catch(error => console.error('Error loading documentation:', error));
}
function closeDocumentation() {
    const modal = document.getElementById('documentationModal');
    modal.style.display = 'none';
}
window.onclick = function(event) {
    const modal = document.getElementById('documentationModal');
    if (event.target === modal) {
        modal.style.display = 'none';
    }
}
