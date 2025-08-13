# AthenAI Model Training Notebooks

This directory contains Jupyter notebooks used for training and experimenting with our AI models. These notebooks can be run in Google Colab or locally.

## Directory Structure

```
notebooks/
├── model_training/           # Production model training notebooks
│   ├── train_exercise_classifier.ipynb     # Main model training pipeline
│   └── model_evaluation.ipynb              # Model evaluation and metrics
└── experiments/             # Experimental notebooks and prototypes
    └── feature_exploration.ipynb           # Data exploration and feature testing
```

## Setup Instructions

1. Open notebooks in Google Colab
2. Mount your Google Drive
3. Clone this repository
4. Install required dependencies:
   ```python
   !pip install transformers datasets torch numpy pandas scikit-learn wandb
   ```

## Workflow

1. Develop and experiment in `experiments/` directory
2. Once approach is validated, move to production training in `model_training/`
3. Track experiments with Weights & Biases
4. Save trained models to Hugging Face Hub

## Contributing

- Keep notebooks clean and well-documented
- Include markdown cells explaining each major step
- Run all cells before committing to ensure clean outputs
- Use relative paths for data loading
