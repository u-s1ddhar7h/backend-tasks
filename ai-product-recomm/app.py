from flask import Flask, jsonify
import pandas as pd
from joblib import load
import numpy as np

app = Flask(__name__)

# Load model and data
model = load('model.joblib')
user_product_matrix = pd.read_pickle('user_product_matrix.pkl')

@app.route('/recommend/<string:user_id>', methods=['GET'])
def get_recommendations(user_id):
    try:
        # Get user vector
        user_vector = user_product_matrix.loc[user_id].values.reshape(1, -1)
        
        # Find similar users
        _, indices = model.kneighbors(user_vector)
        
        # Get products from similar users
        similar_users = user_product_matrix.iloc[indices[0]]
        recommendations = similar_users.mean(axis=0).sort_values(ascending=False)
        
        # Filter out already rated products
        already_rated = user_product_matrix.loc[user_id][user_product_matrix.loc[user_id] > 0].index
        recommendations = recommendations.drop(already_rated, errors='ignore')
        
        return jsonify({
            'user_id': user_id,
            'recommendations': recommendations.head(5).index.tolist()
        })
    except KeyError:
        return jsonify({'error': 'User not found'}), 404

if __name__ == '__main__':
    app.run(debug=True)