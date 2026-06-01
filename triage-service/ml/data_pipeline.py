import pandas as pd
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.impute import SimpleImputer
import os
from typing import Tuple

class DataLoader:
    def __init__(self, file_path: str):
        self.file_path = file_path

    def load_data(self) -> pd.DataFrame:
        if not os.path.exists(self.file_path):
            raise FileNotFoundError(f"Dataset not found at {self.file_path}")
        print(f"Loading raw data from {self.file_path}...")
        return pd.read_csv(self.file_path)

class DataTransformer:
    def __init__(self):
        self.imputer = SimpleImputer(strategy='median')
        self.scaler = StandardScaler()

    def _rename_and_filter_columns(self, df: pd.DataFrame) -> pd.DataFrame:
        column_mapping = {
            'age': 'age',
            'd1_heartrate_max': 'heart_rate',
            'd1_sysbp_max': 'blood_pressure_systolic',
            'temp_apache': 'temperature'
        }
        
        essential_cols = list(column_mapping.keys()) + ['hospital_death', 'apache_4a_icu_death_prob']
        df = df[essential_cols].copy()
        df.rename(columns=column_mapping, inplace=True)
        return df

    def _generate_target_variable(self, df: pd.DataFrame) -> pd.DataFrame:
        """
        Maps ICU survival data to our Triage Acuity System (2: Emergent, 3: Urgent, 4: Less Urgent)
        """
        conditions = [
            (df['hospital_death'] == 1) | (df['apache_4a_icu_death_prob'] > 0.15),
            (df['apache_4a_icu_death_prob'] > 0.05) & (df['apache_4a_icu_death_prob'] <= 0.15)
        ]
        choices = [2, 3]
        
        df['acuity_score'] = np.select(conditions, choices, default=4)
        df.drop(columns=['hospital_death', 'apache_4a_icu_death_prob'], inplace=True)
        return df

    def process(self, df: pd.DataFrame, is_training: bool = True) -> Tuple[np.ndarray, pd.Series]:
        print("Cleaning and transforming data...")
        df = self._rename_and_filter_columns(df)
        df = self._generate_target_variable(df)

        X = df.drop(columns=['acuity_score'])
        y = df['acuity_score']

        print("Engineering new medical features (Shock Index & Elderly Risk)...")

        X['shock_index'] = X['heart_rate'] / X['blood_pressure_systolic'].replace(0, 1)
        X['high_risk_elderly'] = ((X['age'] > 65) & ((X['heart_rate'] > 110) | (X['temperature'] > 38.5))).astype(int)

        print("Handling missing values (Imputation)...")
        if is_training:
            X_imputed = self.imputer.fit_transform(X)
        else:
            X_imputed = self.imputer.transform(X)

        print("Calibrating and Scaling features...")
        if is_training:
            X_scaled = self.scaler.fit_transform(X_imputed)
        else:
            X_scaled = self.scaler.transform(X_imputed)

        return X_scaled, y

    def get_preprocessors(self) -> dict:
        return {
            'imputer': self.imputer,
            'scaler': self.scaler
        }