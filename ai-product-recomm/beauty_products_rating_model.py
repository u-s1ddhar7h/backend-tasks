import pandas as pd
from sklearn.neighbors import NearestNeighbors
from joblib import dump

def train_model():
    # Load sample data
    df = pd.read_csv('./dataset/ratings_Beauty.csv')
    df = df.dropna()
    df = df.head(5000)
    
    # Create user-product matrix
    user_product_matrix = df.pivot(index='UserId', columns='ProductId', values='Rating').fillna(0)
    
    # Train KNN model
    model = NearestNeighbors(n_neighbors=3, metric='cosine')
    model.fit(user_product_matrix)
    
    # Save model and matrix
    dump(model, 'model.joblib')
    user_product_matrix.to_pickle('user_product_matrix.pkl')
    return model, user_product_matrix

if __name__ == '__main__':
    train_model()