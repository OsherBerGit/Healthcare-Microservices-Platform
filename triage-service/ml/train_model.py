import os
import joblib
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from data_pipeline import DataLoader, DataTransformer

class ModelTrainer:
    def __init__(self, n_estimators: int = 150, max_depth: int = 12):
        self.model = RandomForestClassifier(
            n_estimators=n_estimators,
            max_depth=max_depth,
            class_weight='balanced',
            random_state=42
        )

    def train(self, X_train, y_train):
        print("Training Random Forest model on ICU data...")
        self.model.fit(X_train, y_train)

    def evaluate(self, X_test, y_test):
        print("Evaluating model performance...")
        predictions = self.model.predict(X_test)
        print("\n--- Model Evaluation Report ---")
        print(classification_report(y_test, predictions))

    def save_artifacts(self, preprocessors: dict, filepath: str):
        artifacts = {
            'model': self.model,
            'imputer': preprocessors['imputer'],
            'scaler': preprocessors['scaler']
        }
        os.makedirs(os.path.dirname(filepath), exist_ok=True)
        joblib.dump(artifacts, filepath)
        print(f"\nArtifacts (Model, Imputer, Scaler) successfully saved to:\n{filepath}")

if __name__ == "__main__":
    script_dir = os.path.dirname(os.path.abspath(__file__))
    data_path = os.path.join(script_dir, "..", "data", "raw", "icu_patients.csv")
    artifacts_save_path = os.path.join(script_dir, "..", "data", "models", "triage_artifacts.pkl")

    loader = DataLoader(data_path)
    raw_df = loader.load_data()

    transformer = DataTransformer()
    X, y = transformer.process(raw_df, is_training=True)

    print("Splitting dataset into training and testing sets (80/20)...")
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42, stratify=y)

    trainer = ModelTrainer()
    trainer.train(X_train, y_train)
    trainer.evaluate(X_test, y_test)

    preprocessors = transformer.get_preprocessors()
    trainer.save_artifacts(preprocessors, artifacts_save_path)